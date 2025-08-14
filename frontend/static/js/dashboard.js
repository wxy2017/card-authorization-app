// 控制台页面JavaScript

// 加载用户信息
async function loadUserInfo() {
    try {
        const response = await fetch('/api/profile', {
            headers: getAuthHeaders()
        });
        
        if (response.ok) {
            const user = await response.json();
            document.getElementById('showUserNickname').textContent = user.nickname || user.username;
            document.getElementById('showEmail').textContent = user.email;
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

// 加载邮件修改元素
async function loadEmailModal() {
    // 获取元素
    const emailElement = document.getElementById('showEmail');
    const modal = document.getElementById('emailModal');
    const newEmailInput = document.getElementById('newEmail');
    const cancelBtn = document.getElementById('cancelBtn');
    const saveBtn = document.getElementById('saveBtn');
    const emailError = document.getElementById('emailError');

    // 当前邮箱地址
    let currentEmail = emailElement.textContent.trim();

    // 点击邮箱显示弹框
    emailElement.addEventListener('click', () => {
        newEmailInput.value = currentEmail;
        emailError.style.display = 'none';
        modal.style.display = 'flex';
        newEmailInput.focus();
    });

    // 关闭弹框
    cancelBtn.addEventListener('click', () => {
        modal.style.display = 'none';
    });

    // 点击弹框外部关闭
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    });

    // 验证邮箱格式
    function validateEmail(email) {
        const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return regex.test(email);
    }

    // 显示错误信息
    function showError(message) {
        emailError.textContent = message;
        emailError.style.display = 'block';
    }

    // 保存修改并发送到后台
    saveBtn.addEventListener('click', async () => {
        const newEmail = newEmailInput.value.trim();

        // 验证邮箱
        if (!newEmail) {
            showError('请输入邮箱地址');
            return;
        }

        if (!validateEmail(newEmail)) {
            showError('请输入有效的邮箱地址');
            return;
        }

        if (newEmail === currentEmail) {
            showError('新邮箱与当前邮箱相同');
            return;
        }

        try {
            var loginUser = JSON.parse(localStorage.getItem('user') || '{}');
            loginUser.email = newEmail
            // 发送到后台
            const response = await fetch(`/api/user/${loginUser.id}/update`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({ user: loginUser })
            });
            const result = await response.json();
            if (response.ok) {
                // 更新页面显示
                currentEmail = newEmail;
                emailElement.textContent = newEmail;
                modal.style.display = 'none';
                alert('邮箱地址更新成功');
            } else {
                showError(result.message || '更新失败，请稍后重试');
            }
        } catch (error) {
            console.error('更新邮箱错误:', error);
            showError('网络错误，请稍后重试');
        }
    });

    // 按ESC键关闭弹框
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modal.style.display === 'flex') {
            modal.style.display = 'none';
        }
    });

    // 按Enter键保存
    newEmailInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            saveBtn.click();
        }
    });
}

// 加载昵称修改元素
async function loadNikNameModal() {
    // 获取元素
    const nikNameElement = document.getElementById('showUserNickname');
    const modal = document.getElementById('nikNameModal');
    const newNikNameInput = document.getElementById('newNikName');
    const cancelBtn = document.getElementById('nikNameCancelBtn');
    const saveBtn = document.getElementById('nikNameSaveBtn');
    const nikNameError = document.getElementById('nikNameError');

    // 当前邮箱地址
    let currentNikName = nikNameElement.textContent.trim();

    // 点击昵称显示弹框
    nikNameElement.addEventListener('click', () => {
        newNikNameInput.value = currentNikName;
        nikNameError.style.display = 'none';
        modal.style.display = 'flex';
        newNikNameInput.focus();
    });

    // 关闭弹框
    cancelBtn.addEventListener('click', () => {
        modal.style.display = 'none';
    });

    // 点击弹框外部关闭
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    });


    // 显示错误信息
    function showError(message) {
        emailError.textContent = message;
        emailError.style.display = 'block';
    }

    // 保存修改并发送到后台
    saveBtn.addEventListener('click', async () => {
        const newNikName = newNikNameInput.value.trim();

        // 验证邮箱
        if (!newNikName) {
            showError('请输入昵称');
            return;
        }

        try {
            var loginUser = JSON.parse(localStorage.getItem('user') || '{}');
            loginUser.nickname = newNikName
            // 发送到后台
            const response = await fetch(`/api/user/${loginUser.id}/update`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({ user: loginUser })
            });
            const result = await response.json();
            if (response.ok) {
                // 更新页面显示
                currentNikName = newNikName;
                nikNameElement.textContent = newNikName;
                modal.style.display = 'none';
                alert('昵称更新成功');
            } else {
                showError(result.message || '更新失败，请稍后重试');
            }
        } catch (error) {
            console.error('更新昵称错误:', error);
            showError('网络错误，请稍后重试');
        }
    });

    // 按ESC键关闭弹框
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modal.style.display === 'flex') {
            modal.style.display = 'none';
        }
    });

    // 按Enter键保存
    newEmailInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            saveBtn.click();
        }
    });
}

// 页面加载
document.addEventListener('DOMContentLoaded', () => {
    loadUserInfo();
    loadStats();
    loadRecentActivity();
    loadEmailModal();
    loadNikNameModal();
});
