/**
 * TOPBAR MODULE (Legacy)
 * Handles global navigation and user session display.
 */
import { API } from '../../common/lib/api.js';

export const Topbar = {
    init() {
        const logoutBtn = document.getElementById('btn-logout');
        if (logoutBtn) {
            logoutBtn.onclick = () => {
                // Clear Cookie
                document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                API.clearSession();
                window.location.reload();
            };
        }

        const adminBtn = document.getElementById('btn-admin');
        if (adminBtn) {
            // Check role?
            // In legacy, we might just hide it via CSS if not admin, or check API.
            // For now, simple redirect.
            adminBtn.onclick = () => {
                window.location.href = '/admin/';
            };
        }

        // Initialize Dictionary Selector
        this.initDictSelector();

        // Maybe fetch user info to show avatar?
        this.updateUser();
    },

    async initDictSelector() {
        try {
            // In a real implementation we would render a <select> into a specific container
            // but the current HTML structure doesn't have a placeholder for it yet.
            // Requirement says "Add in Topbar Dropdown".
            // Let's dynamically inject it into .top-nav

            let container = document.getElementById('dict-nav-container');
            if (!container) {
                container = document.createElement('div');
                container.id = 'dict-nav-container';
                container.style.marginRight = '10px';

                // Insert before Status
                const nav = document.querySelector('.top-nav');
                const status = document.getElementById('top-status');
                if (nav && status) {
                    nav.insertBefore(container, status);
                }
            }

            const dicts = await API.get('/dictionaries');
            const savedID = localStorage.getItem('vq_dict_id');

            // Find active ones
            const activeDicts = dicts.filter(d => d.is_active);

            if (activeDicts.length === 0) {
                container.innerHTML = '<span style="color:#666; font-size:0.8rem;">NO DICT</span>';
                return;
            }

            // Build Select
            let html = `<select id="dict-select" class="nav-select" style="background:rgba(0,0,0,0.5); color:var(--neon-cyan); border:1px solid var(--neon-cyan); padding:5px; border-radius:4px; margin-right:15px;">`;
            html += `<option value="">ALL DICTIONARIES</option>`;

            activeDicts.forEach(d => {
                const selected = (savedID && parseInt(savedID) === d.id) ? 'selected' : '';
                html += `<option value="${d.id}" ${selected}>${d.name}</option>`;
            });
            html += `</select>`;

            container.innerHTML = html;

            // Bind Change
            const selectEl = document.getElementById('dict-select');
            selectEl.onchange = (e) => {
                const val = e.target.value;
                if (val) {
                    localStorage.setItem('vq_dict_id', val);
                } else {
                    localStorage.removeItem('vq_dict_id');
                }
                // Reload to apply changes cleanly to game engine
                window.location.reload();
            };

        } catch (e) {
            console.error("Failed to load dictionaries", e);
        }
    },

    async updateUser() {
        try {
            // Placeholder
        } catch (e) { }
    }
};
