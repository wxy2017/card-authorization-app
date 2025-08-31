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
    fetch(`/api/users/friends/search?q=${encodeURIComponent(inputValue)}`, {
        headers: getAuthHeaders()
    })
    .then(response => response.json())
    .then(data => {
        if(data.list && data.list.length > 0){
            displayFriends(data.list);
        }
    })
    .catch(error => {
        console.error('搜索用户失败:', error);
    });

}

// 获取好友的邀请状态（）
function getFriendStatusText(status) {
    switch (status) {
        case 'pending':
            return '已邀请';
        case 'accepted':
            return '已接受';
        case 'rejected':
            return '已拒绝';
        default:
            return '可邀请';
    }
}

// 获取好友的操作按钮·
function getFriendActions(item) {
    // 如果好友的邀请状态为空，则可以显示邀请按钮·
    if (item && item.invited === 'default') {
        return `
            <div class="card-actions">
                <button class="btn btn-primary" onclick="sendFriendRequest(${item.user.id})">邀请</button>
            </div>
        `;
    }else{
        //邀请按钮不可点击
        return `
            <div class="card-actions">
                <button class="btn btn-secondary" disabled>已邀请</button>
            </div>
        `;
    }
}

// 显示查询结果
function displayFriends(list) {
    const containerModel = document.getElementById('friendsModal');
     const container = document.getElementById('friendsSearchResults');

    if (list.length === 0) {
        container.innerHTML = '<div class="card"><p class="text-muted">暂无道友</p></div>';
        return;
    }
    containerModel.style.display = 'block';
    container.innerHTML = list.map(item => `
    <div class="card">
        <div class="card-header">
            <h3 class="card-title">
                ${item.user.nickname || item.user.username} 📮 <small>${item.user.email}</small>
            </h3>
            <span class="card-status status-${item.invited}">${getFriendStatusText(item.invited)}</span>
        </div>
        ${getFriendActions(item)}
    </div>
    `).join('');    
}

// 关闭模态框
function closeFriendsModal() {
    const m = document.getElementById('friendsModal');
    if (m) m.style.display = 'none';
}

// 加载我的道友
async function loadMyFriends() {
    // 这里可以添加加载最近活动的逻辑
    const myFriendsInfoElement = document.getElementById('myFriendsInfo');
    try {
        const response = await fetch('/api/users/friends', {
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

// 加载我邀请的道友
async function loadMyInviteFriends() {
    const myInviteFriendsElement = document.getElementById('myInviteFriends');
    try {
        const response = await fetch('/api/users/friends/myInvite/list', {
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
                myInviteFriendsElement.appendChild(friendElement);
            });
        } else {
            // 暂时显示暂无活动
            friendElement.textContent = '暂无邀请，赶快去邀请道友吧！';
            myInviteFriendsElement.appendChild(friendElement);
        }
    } catch (error) {
        console.error('加载用户信息失败:', error);
    }


}

// 加载邀请我的道友
async function loadInviteMyFriends() {
    const myInviteFriendsElement = document.getElementById('inviteMyFriends');
    try {
        const response = await fetch('/api/users/friends/inviteMy/list', {
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
                myInviteFriendsElement.appendChild(friendElement);
            });
        } else {
            // 暂时显示暂无活动
            friendElement.textContent = '暂无邀请';
            myInviteFriendsElement.appendChild(friendElement);
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

    // 加载我的道友
    loadMyFriends()

    // 加载邀请我的道友
    loadInviteMyFriends()

    // 加载我邀请的道友
    loadMyInviteFriends()

});
