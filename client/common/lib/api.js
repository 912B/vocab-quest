/**
 * API MODULE
 * Centralized fetch wrapper for Vocab Quest.
 */

const API_BASE = '/api';

class ApiClient {
    constructor() {
        this.userId = localStorage.getItem('vq_user_id');
        this.role = localStorage.getItem('vq_role');
    }

    setSession(id, role) {
        this.userId = id;
        this.role = role;
        localStorage.setItem('vq_user_id', id);
        localStorage.setItem('vq_role', role);
    }

    clearSession() {
        this.userId = null;
        this.role = null;
        localStorage.removeItem('vq_user_id');
        localStorage.removeItem('vq_role');
    }

    logout() {
        this.clearSession();
    }

    getHeaders() {
        return {
            'Content-Type': 'application/json',
            // 'Authorization': ... (if we had tokens)
        };
    }

    async get(endpoint, params = {}) {
        const url = new URL(API_BASE + endpoint, window.location.origin);
        // Inject User ID automatically if needed
        if (this.userId) url.searchParams.append('user_id', this.userId);

        Object.keys(params).forEach(key => url.searchParams.append(key, params[key]));

        const res = await fetch(url.toString(), {
            headers: this.getHeaders(),
            credentials: 'include' // CRITICAL: Send Cookies
        });

        if (!res.ok) throw new Error(`API Error: ${res.status}`);
        return await res.json();
    }

    async post(endpoint, body) {
        const res = await fetch(API_BASE + endpoint, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify(body),
            credentials: 'include' // CRITICAL: Send Cookies
        });

        if (!res.ok) throw new Error(`API Error: ${res.status}`);
        return await res.json();
    }
}

export const API = new ApiClient();
