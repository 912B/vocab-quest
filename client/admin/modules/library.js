import { API } from '../../common/lib/api.js';

/**
 * LIBRARY MODULE
 * Unified Dictionary & Word Management
 */

let state = {
    dictionaries: [],
    activeDictId: null,
    words: []
};

export async function render(container) {
    container.innerHTML = `
        <div style="display: grid; grid-template-columns: 250px 1fr; gap: 20px; height: 100%;">
            <!-- Left: Dict List -->
            <div id="lib-sidebar" style="border-right: 1px solid rgba(255,255,255,0.1); padding-right: 10px; display:flex; flex-direction:column;">
                <h3 style="color:var(--text-dim); font-size:0.9rem; margin-bottom:10px;">LIBRARIES</h3>
                <div id="dict-list" style="flex:1; overflow-y:auto;">LOADING...</div>
                <button class="btn-cyber btn-cyber-secondary" style="margin-top:10px;" onclick="window.createDict()">+ NEW LIB</button>
            </div>
            
            <!-- Right: Content -->
            <div id="lib-content" style="display:flex; flex-direction:column;">
                <div style="display:flex; justify-content:center; align-items:center; height:100%; color:var(--text-dim);">
                    SELECT A LIBRARY TO MANAGE CONTENT
                </div>
            </div>
        </div>
    `;

    // Bind global actions
    window.selectLibrary = (id) => loadWords(id);
    window.createDict = () => {
        const existing = document.getElementById('dict-create-modal');
        if (existing) existing.remove();

        const modal = document.createElement('div');
        modal.id = 'dict-create-modal';
        modal.className = 'full-screen';
        modal.style.position = 'fixed';
        modal.style.top = '0';
        modal.style.left = '0';
        modal.style.background = 'rgba(0,0,0,0.8)';
        modal.style.zIndex = '200';

        modal.innerHTML = `
            <div class="cyber-card" style="width:95%; max-width:500px;">
                <h3>NEW LIBRARY</h3>
                
                <div style="margin-bottom:15px">
                    <label style="color:var(--text-dim);font-size:0.8rem">LIBRARY NAME</label>
                    <input type="text" id="new-dict-name" class="cyber-input" placeholder="e.g. TOEFL Core">
                </div>
                
                <div style="margin-bottom:15px">
                    <label style="color:var(--text-dim);font-size:0.8rem">DESCRIPTION</label>
                    <textarea id="new-dict-desc" class="cyber-input" rows="3" placeholder="Brief description..."></textarea>
                </div>

                <div style="margin-top:20px; display:flex; gap:10px;">
                        <button id="btn-create-dict" class="btn-cyber">CREATE</button>
                        <button class="btn-cyber btn-cyber-secondary" onclick="document.getElementById('dict-create-modal').remove()">CANCEL</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);

        modal.querySelector('#btn-create-dict').onclick = async () => {
            const name = document.getElementById('new-dict-name').value;
            const desc = document.getElementById('new-dict-desc').value;

            if (!name) return alert("Name is required");

            try {
                await API.post('/dictionaries', {
                    name: name,
                    description: desc,
                    is_active: false
                });
                alert("Library Created.");
                modal.remove();
                loadDictionaries();
            } catch (e) {
                alert("Error: " + e.message);
            }
        };
    };
    window.editWord = (w) => window.openModal(w); // Reuses global modal from index.html
    window.openAddWord = () => window.openModal(null);
    window.deleteWord = deleteWord;

    // Dictionary Management
    window.editDict = (id) => {
        const dict = state.dictionaries.find(d => d.id === id);
        if (!dict) return;

        // Reuse create modal logic or create a dedicated edit one. 
        // For simplicity, we reuse the create structure but populate it.
        const existing = document.getElementById('dict-create-modal');
        if (existing) existing.remove();

        const modal = document.createElement('div');
        modal.id = 'dict-create-modal';
        modal.className = 'full-screen';
        modal.style.position = 'fixed';
        modal.style.top = '0';
        modal.style.left = '0';
        modal.style.background = 'rgba(0,0,0,0.8)';
        modal.style.zIndex = '200';

        modal.innerHTML = `
            <div class="cyber-card" style="width:95%; max-width:500px;">
                <h3>${dict ? 'EDIT LIBRARY' : 'NEW LIBRARY'}</h3>
                
                <div style="margin-bottom:15px">
                    <label style="color:var(--text-dim);font-size:0.8rem">LIBRARY NAME</label>
                    <input type="text" id="new-dict-name" class="cyber-input" value="${dict.name}" placeholder="e.g. TOEFL Core">
                </div>
                
                <div style="margin-bottom:15px">
                    <label style="color:var(--text-dim);font-size:0.8rem">DESCRIPTION</label>
                    <textarea id="new-dict-desc" class="cyber-input" rows="3" placeholder="Brief description...">${dict.description || ''}</textarea>
                </div>

                <div style="margin-top:20px; display:flex; gap:10px;">
                        <button id="btn-create-dict" class="btn-cyber">SAVE</button>
                        <button class="btn-cyber btn-cyber-secondary" onclick="document.getElementById('dict-create-modal').remove()">CANCEL</button>
                </div>
            </div>
        `;
        document.body.appendChild(modal);

        modal.querySelector('#btn-create-dict').onclick = async () => {
            const name = document.getElementById('new-dict-name').value;
            const desc = document.getElementById('new-dict-desc').value;

            if (!name) return alert("Name is required");

            try {
                const payload = {
                    id: dict.id,
                    name: name,
                    description: desc,
                    is_active: dict.is_active
                };

                const res = await fetch('/api/dictionaries', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) throw new Error(await res.text());

                alert("Library Updated.");
                modal.remove();
                loadDictionaries();
            } catch (e) {
                alert("Error: " + e.message);
            }
        };
    };

    window.deleteDict = async (id) => {
        if (!confirm("Are you sure? This will delete ALL words in this library!")) return;
        try {
            const res = await fetch(`/api/dictionaries?id=${id}`, { method: 'DELETE' });
            if (!res.ok) throw new Error(await res.text());
            alert("Library Deleted.");
            // If active was deleted, clear active
            if (state.activeDictId === id) {
                state.activeDictId = null;
                document.getElementById('lib-content').innerHTML = `
                    <div style="display:flex; justify-content:center; align-items:center; height:100%; color:var(--text-dim);">
                        SELECT A LIBRARY TO MANAGE CONTENT
                    </div>`;
            }
            loadDictionaries();
        } catch (e) {
            alert("Delete Failed: " + e.message);
        }
    };

    loadDictionaries();
}

async function loadDictionaries() {
    try {
        state.dictionaries = await API.get('/dictionaries');
        renderDictList();
    } catch (e) {
        document.getElementById('dict-list').innerHTML = `<div style="color:red; font-size:0.8rem; padding:10px;">${e.message}</div>`;
    }
}

function renderDictList() {
    const el = document.getElementById('dict-list');
    el.innerHTML = state.dictionaries.map(d => `
        <div class="nav-item ${state.activeDictId === d.id ? 'active' : ''}" 
             onclick="window.selectLibrary(${d.id})"
             style="margin-bottom:5px; border-radius:4px; border:1px solid ${d.is_active ? 'var(--neon-green)' : 'transparent'}; display:flex; justify-content:space-between; align-items:center; padding-right:5px;">
            <div>
                <div style="font-weight:bold;">${d.name}</div>
                <div style="font-size:0.7rem; opacity:0.7">${d.words_count || 0} WORDS</div>
            </div>
            <div style="display:flex; gap:5px;">
                 <button style="background:none; border:none; cursor:pointer; font-size:1rem; padding:2px; opacity:0.7; transition:opacity 0.2s;" 
                    onclick="event.stopPropagation(); window.editDict(${d.id})" title="Edit">✏️</button>
                 <button style="background:none; border:none; cursor:pointer; font-size:1rem; padding:2px; color:var(--neon-red); opacity:0.7; transition:opacity 0.2s;" 
                    onclick="event.stopPropagation(); window.deleteDict(${d.id})" title="Delete">🗑️</button>
            </div>
        </div>
    `).join('');
}

async function loadWords(dictId) {
    state.activeDictId = dictId;
    renderDictList(); // Update active state

    const content = document.getElementById('lib-content');
    content.innerHTML = "LOADING WORDS...";

    // Set global currentDictId for the modal logic in admin.js to work
    // Ideally we shouldn't rely on globals, but reusing admin.js modal logic requires it or refactoring.
    // Let's rely on admin.js having 'currentDictId' variable exposed? 
    // No, admin.js defines it locally.
    // **FIX**: We need to export this ID or handle save logic HERE.
    // Let's override the global saveWord function or modify how it works.

    // Better: We implement the save logic right here and expose it to window.saveWord
    // essentially capturing the modal form submit.

    try {
        state.words = await API.get('/dictionaries/words', { id: dictId });
        renderWordTable(content, dictId);
    } catch (e) {
        content.innerHTML = "Subsystem Error";
    }
}

function renderWordTable(container, dictId) {
    const dict = state.dictionaries.find(d => d.id === dictId);

    let html = `
        <div>
                <h2 style="margin:0; font-size:1.5rem; color:#fff">${dict.name}</h2>
                <div style="font-size:0.8rem; color:var(--text-dim)">${dict.description}</div>
            </div>
        <div style="display:flex; gap:10px;">
            <!-- Excel Upload -->
            <input type="file" id="inp-excel" accept=".xlsx" style="display:none" onchange="window.uploadExcel(this)">
                <button class="btn-cyber btn-cyber-secondary" onclick="document.getElementById('inp-excel').click()">📄 UPLOAD EXCEL</button>
                <button class="btn-cyber" onclick="window.openAddWord()">+ ADD ENTRY</button>
        </div>
        </div>

        <div style="flex:1; overflow-y:auto;">
            <table class="word-table" style="width:100%; border-collapse:collapse;">
                <thead style="position:sticky; top:0; background:var(--panel-bg);">
                    <tr style="text-align:left; color:var(--neon-cyan);">
                        <th style="padding:10px">WORD</th>
                        <th>DEFINITION</th>
                        <th>LVL</th>
                        <th>ACT</th>
                    </tr>
                </thead>
                <tbody>
                    ${state.words.map(w => `
                        <tr style="border-bottom:1px solid rgba(255,255,255,0.05);">
                            <td style="padding:10px; font-weight:bold; color:#fff">${w.text}</td>
                            <td>${w.definition}</td>
                            <td><span class="tag-diff">L${w.difficulty}</span></td>
                            <td>
                                <button style="background:none; border:none; cursor:pointer;" onclick='window.editWord(${JSON.stringify(w)})'>✏️</button>
                                <button style="background:none; border:none; cursor:pointer; color:red;" onclick="window.deleteWord(${w.id})">🗑️</button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
    container.innerHTML = html;

    // Override Save Handler from admin.js to work with this module's scope
    window.saveWord = async () => {
        const text = document.getElementById('inp-word').value;
        const def = document.getElementById('inp-def').value;
        const diff = parseInt(document.getElementById('inp-diff').value);

        // We know we are in this module, so we use state.activeDictId
        // And we need to check if we are editing (from admin.js global or hidden field?)
        // The modal in HTML doesn't store ID clearly. 
        // admin.js stored it in `editingWordId`.
        // Let's re-implement `editingWordId` tracking here or access admin.js?
        // JS modules have strict scope.

        // Strategy: We inject the ID into the DOM of the modal when opening.
        // OR we just use a closure here for the current editing word.
        // But `window.saveWord` is called by the HTML button `onclick = "saveWord()"`.
        // So we need to override that global function.

        // Let's assume the modal inputs are `inp - word`, etc.
        // We need to know IF we are editing.
        // Hack: Check if `editingWordId` variable on window exists? No.

        // Let's fix the Modal Open to store the ID in a data attribute on the modal.
    };
}

// Re-implementing Modal Logic locally and overriding window globals
// This ensures this module controls the modal without relying on legacy admin.js vars
window.openModal = (word = null) => {
    const overlay = document.getElementById('modal-overlay');
    const title = document.getElementById('modal-title');
    const inpWord = document.getElementById('inp-word');
    const inpDef = document.getElementById('inp-def');
    const inpDiff = document.getElementById('inp-diff');
    const btnSave = overlay.querySelector('button[onclick="saveWord()"]');

    // Hijack click
    btnSave.onclick = () => performSave(word ? word.id : null);

    if (word) {
        title.textContent = "EDIT ENTRY";
        inpWord.value = word.text;
        inpDef.value = word.definition;
        inpDiff.value = word.difficulty;
    } else {
        title.textContent = "NEW ENTRY";
        inpWord.value = "";
        inpDef.value = "";
        inpDiff.value = 1;
    }
    overlay.classList.remove('hidden');
};

async function performSave(editId) {
    const text = document.getElementById('inp-word').value;
    const def = document.getElementById('inp-def').value;
    const diff = parseInt(document.getElementById('inp-diff').value);

    try {
        const payload = {
            dictionary_id: state.activeDictId,
            text: text,
            definition: def,
            difficulty: diff
        };

        if (editId) {
            payload.id = editId;
            await fetch('/api/words', { method: 'PUT', body: JSON.stringify(payload) }); // TODO: Fix body passing
            // The fetch helper in API.js might need work for PUT.
            // Using raw fetch for safety here or API.post?
            // Main.go handle PUT? 
            // api/words handles POST, PUT, DELETE switch.
            await fetch('/api/words', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
        } else {
            await API.post('/words', payload);
        }

        document.getElementById('modal-overlay').classList.add('hidden');
        loadWords(state.activeDictId); // Refresh

    } catch (e) {
        alert("Save Error: " + e.message);
    }
}

async function deleteWord(id) {
    if (!confirm("Delete?")) return;
    await fetch(`/api/words?id=${id}`, { method: 'DELETE' });
    loadWords(state.activeDictId);
}

// Excel Upload Handler
window.uploadExcel = async (input) => {
    if (!input.files || input.files.length === 0) return;

    const file = input.files[0];
    const formData = new FormData();
    formData.append('dictionary_id', state.activeDictId);
    formData.append('file', file);

    const btn = document.querySelector('button[onclick*="inp-excel"]');
    const originalText = btn.textContent;
    btn.textContent = "UPLOADING...";
    btn.disabled = true;

    try {
        const res = await fetch('/api/dictionaries/import', {
            method: 'POST',
            body: formData
        });

        if (!res.ok) throw new Error(await res.text());

        const data = await res.json();
        alert(`SUCCESS: Imported ${data.count} words.`);
        loadWords(state.activeDictId);
    } catch (e) {
        alert("Upload Failed: " + e.message);
    } finally {
        btn.textContent = originalText;
        btn.disabled = false;
        input.value = ""; // Reset
    }
};
