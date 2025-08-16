// 卡片管理页面JavaScript

let currentCardId = null;

// 切换标签
function showTab(tab) {
    document.getElementById('myCards').style.display = tab === 'my' ? 'block' : 'none';
    document.getElementById('receivedCards').style.display = tab === 'received' ? 'block' : 'none';
    document.getElementById('usedCards').style.display = tab === 'used' ? 'block' : 'none';

    // 更新按钮样式
    document.querySelectorAll('.card button').forEach(btn => {
        btn.classList.remove('btn-primary');
        btn.classList.add('btn-secondary');
    });
    event.target.classList.remove('btn-secondary');
    event.target.classList.add('btn-primary');

    // 加载对应数据
    if (tab === 'my') {
        loadMyCards();
    } else if (tab === 'received') {
        loadReceivedCards();
    }else if (tab === 'used') {
        loadUsedCards();
    }
}

// 加载我创建的卡片
async function loadMyCards() {
    try {
        const response = await fetch('/api/cards', {
            headers: getAuthHeaders()
        });

        if (response.ok) {
            const data = await response.json();
            displayCards(data.cards, 'myCards');
        } else if (response.status === 401) {
            logout();
        }
    } catch (error) {
        console.error('加载卡片失败:', error);
    }
}

// 加载收到的卡片
async function loadReceivedCards() {
    try {
        const response = await fetch('/api/cards/received', {
            headers: getAuthHeaders()
        });

        if (response.ok) {
            const data = await response.json();
            displayCards(data.cards, 'receivedCards');
        } else if (response.status === 401) {
            logout();
        }
    } catch (error) {
        console.error('加载卡片失败:', error);
    }
}

async function loadUsedCards() {
    try {
        const response = await fetch(`/api/cards/used`, {
            method: 'POST',
            headers: getAuthHeaders()
        });

        if (response.ok) {
            const data = await response.json();
            displayCards(data.cards, 'usedCards');
        } else if (response.status === 401) {
            logout();
        }
    } catch (error) {
        console.error('加载卡片失败:', error);
    }
}

// 显示卡片
function displayCards(cards, containerId) {
    const container = document.getElementById(containerId);

    if (cards.length === 0) {
        container.innerHTML = '<div class="card"><p class="text-muted">暂无卡片</p></div>';
        return;
    }

    container.innerHTML = cards.map(card => `
        <div class="card">
            <div class="card-header">
                <h3 class="card-title">${card.title}</h3>
                <span class="card-status status-${card.status}">${getStatusText(card.status)}</span>
            </div>
           <div style="display: flex; justify-content: space-between;">
              <span class="card-description">${card.description}</span>
              ${card.status === 'active' ? `
                <span class="card-description">
                    ${getRemainingTime(card.expires_at)}
                </span>` : ''}
            </div>
            <div class="card-meta">
                <span>创建者：${card.creator.nickname || card.creator.username}</span>
                <span>所属者：${card.owner.nickname || card.owner.username}</span>
                <span>${formatDate(card.updated_at)}</span>
            </div>
            ${getCardActions(card)}
        </div>
    `).join('');
}

// 计算剩余天数
function getRemainingTime(expiresAt) {
    if (!expiresAt) return '';
    const now = new Date();
    const expireDate = new Date(expiresAt);
    const diffMs = expireDate - now;
    if (diffMs <= 0) return '';
    const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    return days > 0 ? `<div style="font-weight: 600; font-size: 12px;">剩余<span style="color: #ef4444">${days}</span>天</div>` :
        `<div style="font-weight: 600; font-size: 12px;">剩余<span style="color: #ef4444">${hours}</span>小时</div>`;
}

// 获取状态文本
function getStatusText(status) {
    const statusMap = {
        'active': '可用',
        'used': '已使用',
        'expired': '已过期'
    };
    return statusMap[status] || status;
}

//删除卡
function deleteCard(cardId) {
    try {
        if (!confirm('确定要删除这张卡片吗？')) {
            return;
        }
        fetch(`/api/cards/${cardId}/delete`, {
            method: 'POST',
            headers: getAuthHeaders()
        })
            .then(response => response.json().then(data => ({ response, data })))
            .then(({ response, data }) => {
                if (response.ok) {
                    alert('卡片删除成功！');
                    loadMyCards();
                } else {
                    alert(data.error || '卡片删除失败');
                }
            })
            .catch(error => {
                alert('网络错误，请重试');
            });
    } catch (error) {
        alert('网络错误，请重试');
    }
}

// 获取卡片操作按钮
function getCardActions(card) {
    var loginUser = JSON.parse(localStorage.getItem('user') || '{}');
    var loginUseName = loginUser.username
    if (card.status !== 'active') {
        if(card.status === 'expired' && card.creator.username === loginUseName && card.owner.username === loginUseName ){
            return `
             <div style="margin-top: 0.1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-primary" style="background-color: red; color: white;" onclick="deleteCard(${card.id})">删除</button>
            </div>`;
        }
        return '';
    }

    if (card.creator.username === loginUseName) {
        if(card.owner.username === loginUseName ){
            // 卡的拥有者和创建者有删除权限
            return `
            <div style="margin-top: 1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-primary" onclick="sendCard(${card.id})">发送</button>
            </div>
             <div style="margin-top: 0.1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-primary" style="background-color: red; color: white;" onclick="deleteCard(${card.id})">删除</button>
            </div>
        `;
        }else {
            return `
            <div style="margin-top: 1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-primary" onclick="sendCard(${card.id})">发送</button>
            </div>
        `;
        }
    } else {
        return `
            <div style="margin-top: 1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-success" onclick="useCard(${card.id})">使用</button>
            </div>
        `;
    }
}

// 发送卡片
function sendCard(cardId) {
    currentCardId = cardId;
    document.getElementById('sendModal').style.display = 'block';
}

// 关闭发送模态框
function closeSendModal() {
    document.getElementById('sendModal').style.display = 'none';
    document.getElementById('toUsername').value = '';
}

// 使用卡片
async function useCard(cardId) {
    if (!confirm('确定要使用这张卡片吗？使用后卡片将自动注销。')) {
        return;
    }
    try {
        const response = await fetch(`/api/cards/${cardId}/use`, {
            method: 'POST',
            headers: getAuthHeaders()
        });

        const data = await response.json();

        if (response.ok) {
            alert('卡片使用成功！');
            loadReceivedCards();
        } else {
            alert(data.error || '使用卡片失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
}

// 发送卡片表单提交
document.getElementById('sendForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const toUsername = document.getElementById('toUsername').value;

    try {
        const response = await fetch(`/api/cards/${currentCardId}/send`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({to_username: toUsername})
        });

        const data = await response.json();

        if (response.ok) {
            alert(data.message);
            closeSendModal();
            loadMyCards();
        } else {
            alert(data.error || '发送卡片失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
});

// 格式化日期
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN');
}

// 下拉选项（用户）
async function loadUserList() {
    try {
        const response = await fetch('/api/users/listUsers', {
            headers: getAuthHeaders()
        });
        if (response.ok) {
            const data = await response.json();
            const select = document.getElementById('toUsername');
            select.innerHTML = '<option value="" disabled selected>请选择用户</option>';
            data.users.forEach(user => {
                const option = document.createElement('option');
                option.value = user.username;
                option.textContent = user.nickname;
                select.appendChild(option);
            });
        } else if (response.status === 401) {
            logout();
        }
    } catch (error) {
        console.error('加载用户列表失败:', error);
    }
}

//自动跳转 “收到的”或“我创建的”
async function loadMyAction() {
    const urlParams = new URLSearchParams(window.location.search);
    const tab = urlParams.get('tab');
    if (tab === 'created') {
        document.getElementById('createdCardsTab').click();
    }else if(tab === 'used' ){
        document.getElementById('usedCardsTab').click();
    }else {
        document.getElementById('receivedCardsTab').click();
    }
}

// 页面加载
document.addEventListener('DOMContentLoaded', () => {
    //默认加载我收到的卡
    loadReceivedCards();
    //加载用户列表
    loadUserList();
    // 点击模态框外部关闭
    document.getElementById('sendModal').addEventListener('click', (e) => {
        if (e.target.id === 'sendModal') {
            closeSendModal();
        }
    });
    //加载card?tab=[created|owner]
    loadMyAction();
});


