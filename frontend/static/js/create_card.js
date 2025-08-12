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
        ...(expiresAt && { expiresAt: new Date(expiresAt).toISOString() })
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
            window.location.href = '/cards';
        } else {
            alert(data.error || '创建卡片失败');
        }
    } catch (error) {
        alert('网络错误，请重试');
    }
});
