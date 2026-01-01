/**
 * INPUT MODULE
 * Handles Physical and Virtual Keyboard interactions.
 */

export const Input = {
    game: null, // Check Result Pointer
    state: null, // Shared State Pointer
    ui: null,   // UI Pointer

    listener: null,

    init(gameInstance, stateRef, uiRef) {
        this.game = gameInstance;
        this.state = stateRef;
        this.ui = uiRef;
        this.setupListeners();
    },

    setupListeners() {
        this.listener = (e) => {
            if (this.state.busy || this.state.waitingForRetry) return;
            if (e.key === 'Backspace') this.handle('BACKSPACE');
            else if (/^[a-zA-Z]$/.test(e.key)) this.handle(e.key);
        };
        document.addEventListener('keydown', this.listener);
    },

    dispose() {
        if (this.listener) {
            document.removeEventListener('keydown', this.listener);
            this.listener = null;
        }
    },

    handle(val) {
        if (this.state.busy || this.state.waitingForRetry) return;



        if (!this.state.currWord) return;

        if (val === 'BACKSPACE') {
            // Remove last input
            let lastFilled = -1;
            for (let i = this.state.input.length - 1; i >= 0; i--) {
                if (this.state.mask.includes(i) && this.state.input[i] !== null) {
                    lastFilled = i;
                    break;
                }
            }
            if (lastFilled !== -1) {
                this.state.input[lastFilled] = null;
                this.ui.updateSlotVisuals();
            }
        } else {
            // Add input
            let firstEmpty = -1;
            for (let i = 0; i < this.state.input.length; i++) {
                if (this.state.mask.includes(i) && this.state.input[i] === null) {
                    firstEmpty = i;
                    break;
                }
            }
            if (firstEmpty !== -1) {
                this.state.input[firstEmpty] = val; // Store exact case
                this.ui.updateSlotVisuals();

                // Check Completion
                if (!this.state.input.includes(null)) this.game.checkResult();
            }
        }
    },

    renderKeyboard() {
        const kb = document.getElementById('virtual-keyboard');
        kb.innerHTML = '';

        const rows = [
            "qwertyuiop",
            "asdfghjkl",
            "zxcvbnm"
        ];

        rows.forEach((rowChars, i) => {
            const rowDiv = document.createElement('div');
            rowDiv.className = 'kb-row';

            rowChars.split('').forEach(char => {
                const key = document.createElement('div');
                key.className = 'key';
                key.textContent = char.toUpperCase(); // Legacy: Uppercase Display
                key.onclick = () => this.handle(char);
                rowDiv.appendChild(key);
            });

            // Add Backspace to last row
            if (i === 2) {
                const bs = document.createElement('div');
                bs.className = 'key delete-key';
                bs.innerHTML = '⌫';
                bs.onclick = () => this.handle('BACKSPACE');
                rowDiv.appendChild(bs);
            }

            kb.appendChild(rowDiv);
        });
    }
};
