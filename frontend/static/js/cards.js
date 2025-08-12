// 卡片管理页面JavaScript

let currentCardId = null;

// 切换标签
function showTab(tab) {
    document.getElementById('myCards').style.display = tab === 'my' ? 'block' : 'none';
    document.getElementById('receivedCards').style.display = tab === 'received' ? 'block' : 'none';
    
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
    } else {
        loadReceivedCards();
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
            <p class="card-description">${card.description}</p>
            <div class="card-meta">
                <span>创建者：${card.creator.nickname || card.creator.username}</span>
                <span>${formatDate(card.created_at)}</span>
            </div>
            ${getCardActions(card)}
        </div>
    `).join('');
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

// 获取卡片操作按钮
function getCardActions(card) {
    if (card.status !== 'active') {
        return '';
    }
    
    if (window.location.pathname.includes('myCards')) {
        return `
            <div style="margin-top: 1rem; display: flex; gap: 0.5rem;">
                <button class="btn btn-primary" onclick="sendCard(${card.id})">发送</button>
            </div>
        `;
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
            body: JSON.stringify({ to_username: toUsername })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            alert('卡片发送成功！');
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

// 页面加载
document.addEventListener('DOMContentLoaded', () => {
    loadMyCards();
    
    // 点击模态框外部关闭
    document.getElementById('sendModal').addEventListener('click', (e) => {
        if (e.target.id === 'sendModal') {
            closeSendModal();
        }
    });
});
