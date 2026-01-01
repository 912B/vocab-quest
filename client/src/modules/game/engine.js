import { API } from '../../../common/lib/api.js';
import { Store } from '../../core/Store.js';
import { AssetManager } from '../../utils/AssetManager.js';
import { Audio } from '../../core/Audio.js';
import { Input } from './input.js';

/**
 * VOCAB QUEST - GAME ENGINE (MODULAR V3.5)
 */

const CONFIG = {
    asteroidTimer: 5.00,
    baseScore: 100
};

const STATE = {
    queue: [],
    currWord: null,
    input: [],
    mask: [],
    busy: false,
    score: 0,
    combo: 0,
    waitingForRetry: false,
    totalOps: 0,
    completedOps: 0
};

// --- UI MODULE ---
const UI = {
    el: {
        defBox: null,
        slotCon: null,
        feedback: null,
        asteroid: null,
        kb: null,
        retryBtn: null,
        stamp: null
    },

    init(container) {
        // Use Scoped Query
        this.el.defBox = container.querySelector('#def-box');
        this.el.slotCon = container.querySelector('#slot-container');
        this.el.feedback = container.querySelector('#feedback');
        this.el.asteroid = container.querySelector('#asteroid');
        this.el.kb = container.querySelector('#virtual-keyboard');
        this.el.stamp = container.querySelector('#perfect-stamp');

        // Create Retry Button holder
        const btnDiv = document.createElement('div');
        btnDiv.id = "retry-container";
        btnDiv.className = "retry-container";

        // Insert logic
        if (this.el.kb) {
            this.el.kb.parentNode.insertBefore(btnDiv, this.el.kb.nextSibling);
        } else {
            container.appendChild(btnDiv);
        }
        this.el.retryBtn = btnDiv;
    },

    updateHUD() {
        // Show Remaining Ops
        const progress = `${STATE.completedOps}/${STATE.totalOps}`;
        const score = `${STATE.score} XP`;

        // Target Global Topbar ID
        const topStatus = document.getElementById('top-status');
        if (topStatus) {
            let comboText = "";
            if (STATE.combo > 1) {
                comboText = `<span style="color:var(--neon-magenta); text-shadow:0 0 10px var(--neon-magenta); margin-left:10px;">HYPER-FLUX x${STATE.combo}</span>`;
            }
            topStatus.innerHTML = `
                <span style="color:var(--text-dim)">MISSION:</span> <b>${progress}</b> 
                &nbsp;|&nbsp; 
                <span style="color:var(--neon-gold)">${score}</span>
                ${comboText}
            `;
            topStatus.style.display = 'block';
        }
    },

    renderSlots() {
        if (!STATE.currWord) return;
        this.el.defBox.textContent = STATE.currWord.definition;
        this.el.slotCon.innerHTML = '';

        for (let i = 0; i < STATE.currWord.text.length; i++) {
            const s = document.createElement('div');
            s.className = 'slot';
            // Visual Style: Prefilled slots look different
            if (!STATE.mask.includes(i)) s.classList.add('prefilled');
            this.el.slotCon.appendChild(s);
        }
        this.updateSlotVisuals();
    },

    updateSlotVisuals() {
        const slots = this.el.slotCon.children;

        // Find next empty slot for cursor effect
        let nextEmpty = -1;
        for (let i = 0; i < STATE.input.length; i++) {
            if (STATE.input[i] === null && STATE.mask.includes(i)) {
                nextEmpty = i;
                break;
            }
        }

        for (let i = 0; i < slots.length; i++) {
            const s = slots[i];
            const val = STATE.input[i];
            s.textContent = val || '';

            // Reset dynamic classes
            s.classList.remove('active', 'filled');

            if (STATE.mask.includes(i)) {
                if (val) s.classList.add('filled');
            } else {
                s.classList.add('prefilled');
            }

            if (i === nextEmpty && !STATE.busy && !STATE.waitingForRetry) s.classList.add('active');
        }
    },

    showFeedback(text, color, shake = false) {
        this.el.feedback.textContent = text;
        this.el.feedback.style.color = color;
        if (color === '#ff3333') { // Error color
            this.el.defBox.style.borderColor = "#ff3333";
        } else {
            this.el.defBox.style.borderColor = "";
        }

        if (shake) {
            document.body.animate([
                { transform: 'translate(0,0)' },
                { transform: 'translate(-5px,0)' },
                { transform: 'translate(5px,0)' },
                { transform: 'translate(0,0)' }
            ], { duration: 300 });
        }
    },

    resetVisuals() {
        this.el.feedback.textContent = "";
        this.el.defBox.style.borderColor = "";
        this.el.defBox.className = 'definition-panel';
        this.el.retryBtn.innerHTML = '';

        const slots = this.el.slotCon.children;
        Array.from(slots).forEach(s => s.classList.remove('wrong', 'correct'));
    },

    startHazard() {
        this.el.asteroid.classList.remove('falling');
        void this.el.asteroid.offsetWidth; // Trigger reflow
        this.el.asteroid.classList.add('falling');
    },

    stopHazard() {
        this.el.asteroid.style.opacity = '0';
    },

    markSlots(status) { // status: 'correct' or 'wrong'
        Array.from(this.el.slotCon.children).forEach(s => s.classList.add(status));
    },

    fillCorrectly(text) {
        const slots = this.el.slotCon.children;
        for (let i = 0; i < text.length; i++) {
            slots[i].textContent = text[i];
        }
    },

    showRetryButton(callback) {
        this.el.retryBtn.innerHTML = `<button class="btn-cyber" id="retry-btn">RETRY CONNECTION</button>`;
        document.getElementById('retry-btn').onclick = callback;
    },

    // --- Scaling Logic ---
    handleResize() {
        const hudHeight = 80;
        const availableHeight = window.innerHeight - hudHeight - 20; // 20px buffer
        const availableWidth = window.innerWidth;

        const content = document.querySelector('.target-display');
        if (!content) return;

        // Reset scale to measure natural size
        content.style.transform = 'none';

        const contentHeight = content.scrollHeight;
        const contentWidth = content.scrollWidth;

        // Calculate Scale ratios
        const scaleH = availableHeight / contentHeight;
        const scaleW = availableWidth / (contentWidth + 20); // Width buffer

        // Use the smaller scale
        let scale = Math.min(scaleH, scaleW);
        if (scale > 1) scale = 1;

        if (scale < 1) {
            content.style.transformOrigin = 'top center';
            content.style.transform = `scale(${scale})`;
        }
    },

    showFloatingReward(text) {
        const el = document.createElement('div');
        el.className = 'float-text';
        el.textContent = text;
        document.body.appendChild(el);
        // Auto-remove after animation
        setTimeout(() => el.remove(), 1500);
    }
};


// --- GAME LOGIC MODULE ---
const Game = {
    resizeListener: null,

    async init(container) {
        console.log("Initializing Game Engine V3 (Modular)...");

        try {
            // UI Init
            UI.init(container);

            // Init Input with refs
            Input.init(this, STATE, UI);
            Input.renderKeyboard();

            // Initial Scale
            setTimeout(() => UI.handleResize(), 100);

            this.resizeListener = () => UI.handleResize();
            window.addEventListener('resize', this.resizeListener);

            // Fetch Session with Dictionary Filter
            const dictID = localStorage.getItem('vq_dict_id');
            const query = dictID ? `?dictionary_id=${dictID}` : '';
            const words = await API.get('/session' + query);

            if (!words || words.length === 0) {
                alert("NO MISSIONS AVAILABLE\nReturning to Command Deck");
                Store.set({ route: 'dashboard' }); // Use Store
                return;
            }
            // SETUP QUEUE (Double Pass)
            STATE.queue = words.map(w => ({ ...w, sessionStage: 0 }));
            STATE.totalOps = STATE.queue.length * 2;
            STATE.completedOps = 0;
            STATE.score = 0;
            STATE.combo = 0;

            this.startLevel();

        } catch (e) {
            console.error("Game Engine Crash:", e);
            const defBox = container.querySelector('#def-box');
            if (defBox) defBox.innerHTML = `<span style="color:red; font-size: 1rem;">SYSTEM FAILURE: ${e.message}</span>`;
        }
    },

    dispose() {
        console.log("Disposing Game Engine...");
        Input.dispose();
        if (this.resizeListener) {
            window.removeEventListener('resize', this.resizeListener);
        }

        // Hide/Clear HUD in Topbar
        const topStatus = document.getElementById('top-status');
        if (topStatus) topStatus.style.display = 'none';

        // Stop Audio? (Optional, maybe keep praise playing)
    },

    startLevel() {
        if (STATE.queue.length === 0) {
            this.showMissionComplete();
            return;
        }

        // Pick Sequential (First in Queue) -> Round Robin Flow
        STATE.currWord = STATE.queue.shift();

        const text = STATE.currWord.text;
        const len = text.length;

        // Difficulty Logic
        const isHard = STATE.currWord.sessionStage >= 1;
        const maskRatio = isHard ? 0.7 : 0.35;

        // Reset Input State
        STATE.input = new Array(len).fill(null);
        STATE.busy = false;
        STATE.waitingForRetry = false;

        // Count only letters for difficulty scaling
        const letterCount = text.replace(/[^a-zA-Z]/g, '').length;
        let hideCount = Math.floor(letterCount * maskRatio);
        hideCount = Math.max(1, hideCount);
        if (hideCount > letterCount) hideCount = letterCount;

        console.log(`Word: ${text}, Stage: ${STATE.currWord.sessionStage}, Hidden: ${hideCount}/${letterCount} (Letters Only)`);

        STATE.mask = this.generateMask(text, hideCount);

        // Pre-fill unmasked slots (Preserve Case)
        for (let i = 0; i < len; i++) {
            if (!STATE.mask.includes(i)) {
                STATE.input[i] = text[i];
            }
        }

        UI.resetVisuals();
        UI.updateHUD();
        UI.renderSlots();
        UI.startHazard();
    },

    generateMask(text, count) {
        const letterIndices = [];
        for (let i = 0; i < text.length; i++) {
            if (/[a-zA-Z]/.test(text[i])) letterIndices.push(i);
        }

        const mask = [];
        while (mask.length < count && letterIndices.length > 0) {
            const r = Math.floor(Math.random() * letterIndices.length);
            mask.push(letterIndices[r]);
            letterIndices.splice(r, 1);
        }
        return mask;
    },

    async checkResult() {
        STATE.busy = true;
        const guess = STATE.input.join('');
        const target = STATE.currWord.text;

        // Case-Insensitive Check
        const isSuccess = guess.toLowerCase() === target.toLowerCase();

        if (isSuccess) {
            // --- SUCCESS ---
            STATE.score += CONFIG.baseScore + (STATE.combo * 10);
            STATE.combo++;

            // Progression
            STATE.currWord.sessionStage = (STATE.currWord.sessionStage || 0) + 1;
            STATE.completedOps++;

            // Re-Queue if < 2
            if (STATE.currWord.sessionStage < 2) {
                STATE.queue.push(STATE.currWord);
            }

            UI.markSlots('correct');

            // Show Floating Reward (Right Side)
            let rewardText = "PERFECT";
            if (STATE.combo > 1) rewardText += ` x${STATE.combo}`;
            UI.showFloatingReward(rewardText);

            // Audio: Speak Word Only
            Audio.speak(target);

            try {
                await API.post('/result', { word_id: STATE.currWord.id, success: true });
            } catch (e) {
                console.error("Result Upload Failed:", e);
            }

            // Wait for tension, then move on.
            setTimeout(() => this.nextLevel(), 2500);

        } else {
            // --- FAILURE ---
            STATE.combo = 0;

            // Remedial: Push back
            if (!STATE.waitingForRetry) {
                STATE.queue.push(STATE.currWord);
            }

            if (UI.el.stamp) UI.el.stamp.classList.remove('visible');

            UI.markSlots('wrong');
            UI.showFeedback(`${target}`, "var(--neon-green)", true);

            // Audio: Speak Correction
            Audio.speak(target);

            try {
                await API.post('/result', { word_id: STATE.currWord.id, success: false });
            } catch (e) {
                console.error("Result Upload Failed:", e);
            }

            // Enter Retry State
            STATE.waitingForRetry = true;
            setTimeout(() => {
                UI.showRetryButton(() => this.resetForRetry());

                // Allow Enter Key for convenience
                const enterHandler = (e) => {
                    if (e.key === 'Enter' && STATE.waitingForRetry) {
                        document.removeEventListener('keydown', enterHandler);
                        this.resetForRetry();
                    }
                };
                document.addEventListener('keydown', enterHandler);
            }, 2000);
        }
    },

    resetForRetry() {
        STATE.waitingForRetry = false;
        STATE.busy = false;

        // Clear ONLY masked slots
        for (let i = 0; i < STATE.input.length; i++) {
            if (STATE.mask.includes(i)) STATE.input[i] = null;
        }

        UI.resetVisuals(); // Clears feedback, retry button, wrong class
        UI.updateSlotVisuals();
    },

    nextLevel() {
        UI.el.asteroid.style.opacity = '1';
        this.startLevel();
    },

    showMissionComplete() {
        const overlay = document.createElement('div');
        overlay.style.cssText = `
            position: fixed; top: 0; left: 0; width: 100%; height: 100%;
            background: rgba(0,0,0,0.9); z-index: 2000;
            display: flex; align-items: center; justify-content: center;
        `;
        overlay.innerHTML = `
            <div class="panel-cyber" style="text-align: center; border: 2px solid var(--neon-gold); box-shadow: 0 0 30px var(--neon-gold);">
                <h2 style="color:var(--neon-gold); font-size: 2rem; margin-bottom: 20px;">MISSION ACCOMPLISHED</h2>
                <div style="margin-bottom: 30px; font-size: 1.2rem; color: #fff;">
                    FINAL SCORE: <b style="color:var(--neon-blue)">${STATE.score}</b> XP
                </div>
                <button class="btn-cyber" id="btn-return" style="font-size: 1.2rem; padding: 15px 40px;">RETURN TO BASE</button>
            </div>
        `;
        document.body.appendChild(overlay);

        document.getElementById('btn-return').onclick = () => {
            overlay.remove();
            Store.set({ route: 'dashboard' });
        };
    }
};

// Export Singleton
export const Engine = Game;
