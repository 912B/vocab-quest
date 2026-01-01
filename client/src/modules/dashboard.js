import { API } from '../../common/lib/api.js';
import { Store } from '../core/Store.js';

export const Dashboard = {
    async init(container) {
        container.innerHTML = `<div class="loading">Loading Mission Data...</div>`;
        try {
            const dictID = localStorage.getItem('vq_dict_id');
            const query = dictID ? `?dictionary_id=${dictID}` : '';
            const stats = await API.get('/stats' + query);
            this.render(container, stats);
        } catch (e) {
            container.innerHTML = `<div class="error">DATA ERROR: ${e.message}</div>`;
        }
    },

    render(container, stats) {
        // Calculate percentages if needed, or just show raw numbers
        // Stats: total_words, mastered_words, learning_words, new_words

        container.innerHTML = `
            <div class="dashboard-panel">
                <h2 style="color:var(--neon-cyan); letter-spacing:2px; margin-bottom:30px;">PILOT STATISTICS</h2>
                
                <div class="stats-grid" style="display:grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap:20px; width:100%; max-width:800px;">
                    
                    <div class="stat-card" style="border:1px solid var(--neon-blue); padding:20px; border-radius:8px; background:rgba(0,0,0,0.3); text-align:center;">
                        <div style="font-size:3rem; font-weight:bold; color:var(--neon-blue);">${stats.total_words}</div>
                        <div style="color:var(--text-dim); font-size:0.8rem; letter-spacing:1px;">TOTAL MISSIONS</div>
                    </div>

                    <div class="stat-card" style="border:1px solid var(--neon-gold); padding:20px; border-radius:8px; background:rgba(0,0,0,0.3); text-align:center;">
                        <div style="font-size:3rem; font-weight:bold; color:var(--neon-gold);">${stats.mastered_words}</div>
                        <div style="color:var(--text-dim); font-size:0.8rem; letter-spacing:1px;">MASTERED</div>
                    </div>

                    <div class="stat-card" style="border:1px solid var(--neon-magenta); padding:20px; border-radius:8px; background:rgba(0,0,0,0.3); text-align:center;">
                        <div style="font-size:3rem; font-weight:bold; color:var(--neon-magenta);">${stats.learning_words}</div>
                        <div style="color:var(--text-dim); font-size:0.8rem; letter-spacing:1px;">IN PROGRESS</div>
                    </div>

                    <div class="stat-card" style="border:1px solid #fff; padding:20px; border-radius:8px; background:rgba(0,0,0,0.3); text-align:center;">
                        <div style="font-size:3rem; font-weight:bold; color:#fff;">${stats.new_words}</div>
                        <div style="color:var(--text-dim); font-size:0.8rem; letter-spacing:1px;">UNEXPLORED</div>
                    </div>

                </div>

                <button class="btn-cyber" style="margin-top:40px;" onclick="location.reload()">RESUME MISSION</button>
            </div>
        `;
    }
};
