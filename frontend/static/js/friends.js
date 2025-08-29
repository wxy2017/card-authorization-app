// 认证相关JavaScript

// 获取认证头
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

// 发送卡片
function searchUser() {
    let inputValue = document.getElementById('searchInput').value;
    if (!inputValue) {
        alert('请输入用户名或邮箱');
        return;
    }
    fetch(`/api/users/search?q=${encodeURIComponent(inputValue)}`, {
        headers: getAuthHeaders()
    })
    .then(response => response.json())
    .then(data => {
        if(data.users && data.users.length > 0){
            for (let user of data.users) {
              alert(`找到用户: ${user.nickname} (${user.email})`);
            }

        }
    })
    .catch(error => {
        console.error('搜索用户失败:', error);
    });

}

// 加载我的道友
async function loadMyFriends() {
    // 这里可以添加加载最近活动的逻辑
    const myFriendsInfoElement = document.getElementById('myFriendsInfo');
    try {
        const response = await fetch('/api/friends', {
            headers: getAuthHeaders()
        });
        const friendElement = document.createElement('div');
        const data = await response.json();
        if (response.ok && data.users && data.users.length > 0) {
            data.users.forEach(user => {
                friendElement.classList.add('card');
                //收到对方发的卡
                friendElement.innerHTML = `
                <h4 style="display: flex; align-items: center;">
                    <span class="gradient-text">${user.nickname}</span>
                    <span>📮<small>${user.email}</small></span>
                </h4>
                `;
                myFriendsInfoElement.appendChild(friendElement);
            });
        } else {
            // 暂时显示暂无活动
            friendElement.textContent = '暂无道友';
            myFriendsInfoElement.appendChild(friendElement);
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
    }


}


// 页面加载时检查认证
document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('token');
    const currentPath = window.location.pathname;

    // 如果已登录，跳转到dashboard
    if (token && (currentPath === '/' || currentPath === '/login' || currentPath === '/register')) {
        window.location.href = '/friends';
    }
    
    // 如果未登录，跳转到登录页
    if (!token && !['/login', '/register'].includes(currentPath)) {
        window.location.href = '/login';
    }

    loadMyFriends()
});
