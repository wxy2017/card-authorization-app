// 创建卡片页面JavaScript

// 创建卡片
document.getElementById('createCardForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const expiresAt = document.getElementById('expiresAt').value;

    const cardData = {
        title,
        description,
        ...(expiresAt && { expires_at: new Date(expiresAt + 'T00:00:00+08:00').toLocaleString('sv-SE', { timeZone: 'Asia/Shanghai' }).replace(' ', 'T') + '+08:00' })
    };

    try {
        const response = await fetch('/api/cards', {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(cardData)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            alert('卡片创建成功！');
            window.location.href = '/cards?tab=created';
        } else {
            alert(data.error || '创建卡片失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
});

// 页面加载完成后设置最小可选日期
document.addEventListener('DOMContentLoaded', function() {
    const expiresAtInput = document.getElementById('expiresAt');
    // 格式化日期为 YYYY-MM-DD 格式（date 输入框要求的标准格式）
    // 设置最小可选日期为今天，确保只能选择今天及之后的日期
    expiresAtInput.min =  new Date().toISOString().split('T')[0];
});
