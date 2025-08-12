// 控制台页面JavaScript

// 加载用户信息
async function loadUserInfo() {
    try {
        const response = await fetch('/api/profile', {
            headers: getAuthHeaders()
        });
        
        if (response.ok) {
            const user = await response.json();
            document.getElementById('userNickname').textContent = user.nickname || user.username;
            localStorage.setItem('user', JSON.stringify(user));
        } else if (response.status === 401) {
            logout();
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
    }
}

// 加载统计信息
async function loadStats() {
    try {
        const [myCardsResponse, receivedCardsResponse] = await Promise.all([
            fetch('/api/cards', { headers: getAuthHeaders() }),
            fetch('/api/cards/received', { headers: getAuthHeaders() })
        ]);
        
        if (myCardsResponse.ok && receivedCardsResponse.ok) {
            const myCards = await myCardsResponse.json();
            const receivedCards = await receivedCardsResponse.json();
            
            document.getElementById('createdCards').textContent = myCards.cards.length;
            document.getElementById('receivedCards').textContent = receivedCards.cards.length;
            
            // 计算已使用的卡片
            const usedCards = [...myCards.cards, ...receivedCards.cards].filter(card => card.status === 'used').length;
            document.getElementById('usedCards').textContent = usedCards;
        }
    } catch (error) {
        console.error('加载统计信息失败:', error);
    }
}

// 加载最近活动
async function loadRecentActivity() {
    // 这里可以添加加载最近活动的逻辑
    // 暂时显示暂无活动
}

// 页面加载
document.addEventListener('DOMContentLoaded', () => {
    loadUserInfo();
    loadStats();
    loadRecentActivity();
});
