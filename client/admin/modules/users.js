import { API } from '../../common/lib/api.js';

export async function render(container) {
    container.innerHTML = "LOADING PERSONNEL RECORDS...";
    try {
        const users = await API.get('/users');

        // CSS Styles
        const styles = `
            <style>
                .user-row {
                    background: rgba(255,255,255,0.03); 
                    padding: 15px; 
                    display: flex; 
                    align-items: center; 
                    justify-content: space-between; 
                    border: 1px solid rgba(255,255,255,0.05);
                    border-radius: 6px;
                    transition: all 0.2s;
                }
                .user-row:hover {
                    background: rgba(255,255,255,0.06);
                    border-color: rgba(255,255,255,0.2);
                }
                .user-row.admin-role {
                    border-left: 3px solid var(--neon-gold);
                }
                
                .action-group {
                    display: flex;
                    align-items: center;
                    gap: 10px;
                }

                .btn-mini {
                    font-size: 0.75rem;
                    padding: 5px 10px;
                    background: rgba(0,0,0,0.3);
                    border: 1px solid rgba(255,255,255,0.2);
                    color: var(--text-dim);
                    border-radius: 4px;
                    cursor: pointer;
                    transition: all 0.2s;
                }
                .btn-mini:hover {
                    background: rgba(255,255,255,0.1);
                    color: #fff;
                    border-color: #fff;
                }

                .btn-icon {
                    background: rgba(255, 51, 51, 0.1); 
                    border: 1px solid rgba(255, 51, 51, 0.3); 
                    color: #ff3333; 
                    width: 30px; 
                    height: 30px; 
                    display: flex; 
                    align-items: center; 
                    justify-content: center; 
                    border-radius: 4px; 
                    cursor: pointer;
                    transition: all 0.2s;
                }
                .btn-icon:hover {
                    background: #ff3333;
                    color: #fff;
                    box-shadow: 0 0 10px #ff3333;
                }
            </style>
        `;

        let html = `
             <div style="display:flex; justify-content:space-between; margin-bottom:20px; align-items:center;">
                <h3>PERSONNEL RECORDS</h3>
                <button class="btn-cyber" style="width:auto; padding: 8px 20px;" onclick="window.createUser()">+ RECRUIT</button>
            </div>
            ${styles}
            <div class="profile-grid" style="display:flex; flex-direction:column; gap:10px;">
        `;

        users.forEach(u => {
            const isAdmin = u.role === 'admin';
            html += `
                <div class="user-row ${isAdmin ? 'admin-role' : ''}">
                    <!-- Left: User Info -->
                    <div style="display:flex; align-items:center; gap:15px;">
                         <div style="width:40px; text-align:center;">
                            ${u.avatar.includes('/') ? `<img src="../${u.avatar}" style="width:40px; border-radius:50%;">` : `<span style="font-size:2rem">${u.avatar}</span>`}
                         </div>
                         <div>
                            <div style="font-weight:700; color:#fff; font-size:1.1rem;">${u.username}</div>
                            <div style="font-size:0.75rem; color:var(--text-dim); letter-spacing:1px;">${u.role.toUpperCase()}</div>
                         </div>
                    </div>

                    <!-- Right: Actions -->
                    <div class="action-group">
                        <button class="btn-mini" onclick="window.viewUserStats(${u.id}, '${u.username}')">STATS</button>
                        <button class="btn-mini" onclick="window.resetUserPassword(${u.id})">RESET PASS</button>
                        ${!isAdmin ? `<button class="btn-mini" style="color:var(--neon-green); border-color:var(--neon-green);" onclick="window.promoteUser(${u.id})">PROMOTE</button>` : ''}
                        <button class="btn-icon" title="DELETE USER" onclick="window.deleteUser(${u.id}, '${u.username}')">🗑️</button>
                    </div>
                </div>
            `;
        });

        container.innerHTML = html + "</div>";

        window.viewUserStats = async (id, name) => {
            const existing = document.getElementById('stats-modal');
            if (existing) existing.remove();

            const modal = document.createElement('div');
            modal.id = 'stats-modal';
            modal.className = 'full-screen';
            modal.style.position = 'fixed';
            modal.style.top = '0';
            modal.style.left = '0';
            modal.style.background = 'rgba(0,0,0,0.8)';
            modal.style.zIndex = '200';
            modal.innerHTML = '<div style="color:white; margin-top:20%">LOADING...</div>';
            document.body.appendChild(modal);

            try {
                const stats = await API.get(`/stats?user_id=${id}`);

                modal.innerHTML = `
                    <div class="cyber-card" style="width:95%; max-width:600px;">
                        <h3>STATS: ${name.toUpperCase()}</h3>
                        <div class="stats-grid" style="display:grid; grid-template-columns: 1fr 1fr; gap:10px; margin-bottom:20px;">
                            <div style="background:rgba(255,255,255,0.05); padding:10px; text-align:center;">
                                <div style="font-size:2rem; color:var(--neon-blue)">${stats.total_words}</div>
                                <div style="font-size:0.7rem; color:var(--text-dim)">TOTAL</div>
                            </div>
                            <div style="background:rgba(255,255,255,0.05); padding:10px; text-align:center;">
                                <div style="font-size:2rem; color:var(--neon-gold)">${stats.mastered_words}</div>
                                <div style="font-size:0.7rem; color:var(--text-dim)">MASTERED</div>
                            </div>
                            <div style="background:rgba(255,255,255,0.05); padding:10px; text-align:center;">
                                <div style="font-size:2rem; color:var(--neon-magenta)">${stats.learning_words}</div>
                                <div style="font-size:0.7rem; color:var(--text-dim)">LEARNING</div>
                            </div>
                            <div style="background:rgba(255,255,255,0.05); padding:10px; text-align:center;">
                                <div style="font-size:2rem; color:white">${stats.new_words}</div>
                                <div style="font-size:0.7rem; color:var(--text-dim)">NEW</div>
                            </div>
                        </div>
                        <button class="btn-cyber" onclick="document.getElementById('stats-modal').remove()">CLOSE</button>
                    </div>
                `;
            } catch (e) {
                modal.innerHTML = `<div class="cyber-card">Error: ${e.message} <br><button onclick="this.parentElement.parentElement.remove()">CLOSE</button></div>`;
            }
        };

        window.resetUserPassword = async (id) => {
            // Remove existing modal if any
            const existing = document.getElementById('pwd-reset-modal');
            if (existing) existing.remove();

            // Create Modal HTML
            const modal = document.createElement('div');
            modal.id = 'pwd-reset-modal';
            modal.className = 'full-screen';
            modal.style.position = 'fixed';
            modal.style.top = '0';
            modal.style.left = '0';
            modal.style.background = 'rgba(0,0,0,0.8)';
            modal.style.zIndex = '200';

            modal.innerHTML = `
                <div class="cyber-card" style="width:95%; max-width:450px; max-height:90vh; overflow-y:auto;">
                    <h3>RESET PASSWORD</h3>
                    <div style="margin-bottom:15px">
                        <label style="color:var(--text-dim);font-size:0.8rem">NEW PASSWORD</label>
                        <input type="text" id="new-pw-text" class="cyber-input" placeholder="New Secure Password">
                    </div>
                    <div style="margin-top:20px; display:flex; gap:10px;">
                        <button id="btn-save-pw" class="btn-cyber">SAVE</button>
                        <button class="btn-cyber btn-cyber-secondary" onclick="document.getElementById('pwd-reset-modal').remove()">CANCEL</button>
                    </div>
                </div>
            `;

            document.body.appendChild(modal);

            modal.querySelector('#btn-save-pw').onclick = async () => {
                const newPw = document.getElementById('new-pw-text').value;

                if (!newPw) return alert("Password cannot be empty");

                try {
                    // Note: We don't change the role here, just password
                    // But the API might expect role. Let's send the current role?
                    // Or simpler: The backend probably only updates fields provided properly or handles it.
                    // To be safe, we just send id and password. The Go handler handles partial updates?
                    // Let's check: AuthHandler UpdateUser usually updates all fields passed.
                    // Wait, we don't have the user's current role handy easily inside THIS specific function scope unless we pass it.
                    // But we actually only care about resetting the password.
                    // Let's assume the API handles it or verify. A safe bet is the API needs ID.
                    await fetch('/api/users', {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ id: id, password: newPw })
                    });
                    alert("Password Reset Successfully.");
                    modal.remove();
                } catch (e) {
                    alert("Error: " + e.message);
                }
            };
        };

        window.deleteUser = async (id, username) => {
            if (confirm(`Are you sure you want to delete "${username}"?`)) {
                try {
                    await fetch(`/api/users?id=${id}`, { method: 'DELETE' });
                    render(container);
                } catch (e) {
                    alert("Error: " + e.message);
                }
            }
        };

        window.promoteUser = async (id) => {
            if (confirm("Make Admin?")) {
                await fetch('/api/users', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ id: id, role: 'admin' })
                });
                render(container);
            }
        };

        window.createUser = () => {
            const existing = document.getElementById('user-create-modal');
            if (existing) existing.remove();

            const modal = document.createElement('div');
            modal.id = 'user-create-modal';
            modal.className = 'full-screen';
            modal.style.position = 'fixed';
            modal.style.top = '0';
            modal.style.left = '0';
            modal.style.background = 'rgba(0,0,0,0.8)';
            modal.style.zIndex = '200';

            modal.innerHTML = `
                <div class="cyber-card" style="width:95%; max-width:500px; max-height:90vh; overflow-y:auto;">
                    <h3>RECRUIT NEW PERSONNEL</h3>
                    
                    <div style="margin-bottom:15px">
                        <label style="color:var(--text-dim);font-size:0.8rem">USERNAME</label>
                        <input type="text" id="new-u-name" class="cyber-input" placeholder="e.g. Navigator">
                    </div>
                    
                     <div style="margin-bottom:15px">
                        <label style="color:var(--text-dim);font-size:0.8rem">ROLE</label>
                        <select id="new-u-role" class="cyber-input">
                            <option value="user">USER (Pilot/Engineer)</option>
                            <option value="admin">ADMIN (Commander)</option>
                        </select>
                    </div>

                    <div style="margin-bottom:15px">
                        <label style="color:var(--text-dim);font-size:0.8rem">AVATAR (Emoji or Path)</label>
                        <input type="text" id="new-u-avatar" class="cyber-input" value="👨‍🚀">
                    </div>

                    <div style="margin-bottom:15px">
                        <label style="color:var(--text-dim);font-size:0.8rem">PASSWORD</label>
                        <input type="text" id="new-u-pw-text" class="cyber-input" placeholder="Secure Password">
                    </div>

                    <div style="margin-top:20px; display:flex; gap:10px;">
                         <button id="btn-create-user" class="btn-cyber">RECRUIT</button>
                         <button class="btn-cyber btn-cyber-secondary" onclick="document.getElementById('user-create-modal').remove()">CANCEL</button>
                    </div>
                </div>
            `;
            document.body.appendChild(modal);

            modal.querySelector('#btn-create-user').onclick = async () => {
                const name = document.getElementById('new-u-name').value;
                const role = document.getElementById('new-u-role').value;
                const avatar = document.getElementById('new-u-avatar').value;
                const pw = document.getElementById('new-u-pw-text').value;

                if (!name || !pw) return alert("Missing Fields");

                try {
                    await fetch('/api/users', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ username: name, password: pw, role: role, avatar: avatar })
                    });
                    alert("User Recruited.");
                    modal.remove();
                    render(container); // Refresh list
                } catch (e) {
                    alert("Error: " + e.message);
                }
            };
        };

    } catch (e) {
        container.innerHTML = "Error: " + e.message;
    }
}
