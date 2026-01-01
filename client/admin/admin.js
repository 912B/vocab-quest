/**
 * ADMIN ROUTER
 * Dynamically loads management modules.
 */

// Global Router
window.router = {
    loadModule: async (name, params = {}) => {
        const content = document.getElementById('content-area');

        // Highlight Sidebar
        document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
        const navEl = document.getElementById(`nav-${name}`);
        if (navEl) navEl.classList.add('active');

        try {
            const module = await import(`./modules/${name}.js?v=2`);
            module.render(content, params);
        } catch (e) {
            content.innerHTML = `<div style="color:red">MODULE LOAD ERROR: ${e.message}</div>`;
        }
    }
};

// Global Modal Utilities
window.closeModal = () => {
    document.getElementById('modal-overlay').classList.add('hidden');
};

// Update Sidebar HTML mapping
document.addEventListener('DOMContentLoaded', () => {
    // Default load
    window.router.loadModule('library');
});
