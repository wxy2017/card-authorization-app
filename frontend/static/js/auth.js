// 认证相关JavaScript

// 检查是否已登录
function checkAuth() {
    const token = localStorage.getItem('token');
    if (token && !window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
        return true;
    }
    return false;
}

// Cookie操作工具函数
function setCookie(name, value, daysToLive) {
    // 计算过期时间
    const date = new Date();
    date.setTime(date.getTime() + (daysToLive * 24 * 60 * 60 * 1000));
    const expires = "expires=" + date.toUTCString();

    // 存储Cookie（包含路径和过期时间）
    document.cookie = `${name}=${encodeURIComponent(value)}; ${expires}; path=/`;
}

function getCookie(name) {
    // 查找指定Cookie
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);

    if (parts.length === 2) {
        return decodeURIComponent(parts.pop().split(';').shift());
    }
    return null;
}

function removeCookie(name) {
    // 通过设置过期时间为过去来删除Cookie
    setCookie(name, '', -1);
}

// 登录
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const rememberMe = document.getElementById('rememberMe')?.checked;
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
            if (rememberMe) {
                // 勾选了记住密码，存储30天
                setCookie('userName', username, 30);
                // 注意：实际项目中必须加密存储密码，这里仅为示例
                setCookie('userPwd', btoa(password), 30);
            } else {
                // 未勾选，清除已存储的Cookie
                removeCookie('userName');
                removeCookie('userPwd');
            }
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

// 获取记住密码后的用户信息
async function loadRememberMeInfo(){
    const savedUser = getCookie('userName');
    const savedPwd = getCookie('userPwd');
    if (savedUser && savedPwd) {
        // 填充表单
        document.getElementById('username').value = savedUser;
        document.getElementById('password').value = atob(savedPwd); // 解码
        document.getElementById('rememberMe').checked = true; // 勾选复选框
    }else{
        // document.getElementById('username').value = "";
        document.getElementById('password').value = "";
        document.getElementById('rememberMe').checked = false;
    }
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

    //页面渲染完毕后执行记住密码信息加载
    if (!token && ['/login'].includes(currentPath)) {
        // 监听页面DOM加载完成事件
        loadRememberMeInfo();
    }
});
