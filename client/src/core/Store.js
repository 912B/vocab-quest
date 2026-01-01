// Minimal Store implementation to support legacy engine.js
export const Store = {
    state: {
        route: 'game',
        user: null
    },

    // Legacy engine uses 'set'
    set(newState) {
        this.state = { ...this.state, ...newState };
        console.log("Store Update:", this.state);

        // If route changes to dashboard, reload for now since we are in legacy mode
        if (newState.route === 'dashboard') {
            window.location.reload();
        }
    },

    get() {
        return this.state;
    }
};
