// 认证相关JavaScript

// 检查是否已登录
function checkAuth() {
    const token = localStorage.getItem('token');
    if (token && !window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
        return true;
    }
    return false;
}

// 登录
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            window.location.href = '/dashboard';
        } else {
            alert(data.error || '登录失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
});

// 注册
document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const nickname = document.getElementById('nickname').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, nickname, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            window.location.href = '/dashboard';
        } else {
            alert(data.error || '注册失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
});

// 退出登录
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
}

// 获取认证头
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

// 页面加载时检查认证
document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('token');
    const currentPath = window.location.pathname;
    
    // 如果已登录，跳转到dashboard
    if (token && (currentPath === '/' || currentPath === '/login' || currentPath === '/register')) {
        window.location.href = '/dashboard';
    }
    
    // 如果未登录，跳转到登录页
    if (!token && !['/login', '/register'].includes(currentPath)) {
        window.location.href = '/login';
    }
});
