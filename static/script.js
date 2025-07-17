async function init() {

    setLoadingState(true);

    try {
        // 发起API请求
        const response = await fetch('http://124.221.182.104:8055/api/index');
        if (!response.ok) {
            throw new Error(`请求失败，状态码: ${response.status}`);
        }

        const data = await response.json();
        console.log('接口数据获取成功:', data);

        const requiredFields = ['推荐阅读', '分类', 'new_in', 'new_up', 'new_msg'];
        requiredFields.forEach(field => {
            if (!data.hasOwnProperty(field)) {
                throw new Error(`接口数据缺少必要字段: ${field}`);
            }
        });

        renderRecommend(data['推荐阅读']);
        renderCategoryGroups(data['分类']);
        renderNewIn(data['new_in']);
        renderNewUp(data['new_up']);
        renderNewMsg(data['new_msg']);

    } catch (error) {
        console.error('数据处理失败:', error);
        // 显示错误状态
        document.querySelectorAll('.loading').forEach(el => {
            el.innerHTML = '<span style="color: #e4393c;">加载失败，请刷新重试</span>';
        });
    } finally {
        setLoadingState(false);
    }
}

// 设置所有区域的加载状态
function setLoadingState(isLoading) {
    const loaders = document.querySelectorAll('.loading');
    loaders.forEach(loader => {
        if (isLoading) {
            loader.innerHTML = '<i class="fa fa-circle-o-notch fa-spin"></i><p>加载中...</p>';
        } else {
            loader.style.display = 'none';
        }
    });
}

// 渲染推荐阅读区域
function renderRecommend(books) {
    const container = document.getElementById('recommend-container');
    container.innerHTML = '';

    if (!books || books.length === 0) {
        container.innerHTML = '<div class="text-center">暂无推荐书籍</div>';
        return;
    }

    books.forEach(book => {
        const bookItem = document.createElement('div');
        bookItem.className = 'book-item';
        bookItem.innerHTML = `
            <div class="book-card">
                <div class="cover-container">
                    <a href="/detail" target="_blank">
                        <div class="thumbnail" style="background-image: url('${book.cover_url || ''}')"></div>
                    </a>
                </div>
                <div class="book-info">
                    <h4 class="book-title">
                        <a href="/detail" target="_blank" title="${book.title || ''}">${book.title || ''}</a>
                    </h4>
                    <div class="book-author">作者：${book.author || '未知作者'}</div>
                    <div class="book-intro">${book.intro || '暂无简介'}</div>
                </div>
            </div>
        `;
        container.appendChild(bookItem);
    });
}


function renderCategoryGroups(categoryData) {
    // 分类容器ID与对应名称的映射（新增仙侠修真）
    const categoryContainers = [
        { id: 'category-xuanhuan', name: '玄幻魔法' },
        { id: 'category-dushi', name: '都市言情' },
        { id: 'category-lishi', name: '历史军事' },
        { id: 'category-wangyou', name: '网游竞技' },
        { id: 'category-kehuan', name: '科幻未来' },
        { id: 'category-xianxia', name: '仙侠修真' } // 新增分类
    ];

    categoryContainers.forEach((container, index) => {
        const start = index * 10; // 每类10条（0-9、10-19...50-59）
        const end = start + 10;
        const items = categoryData.slice(start, end); // 从60条中截取对应分片
        const containerEl = document.getElementById(container.id);
        containerEl.innerHTML = '';
        if (items.length === 0) {
            containerEl.innerHTML = '<li class="list-group-item">暂无数据</li>';
            return;
        }
        // 渲染当前分类的10条数据
        items.forEach(item => {
            const listItem = document.createElement('li');
            listItem.className = 'list-group-item';
            listItem.innerHTML = `
                <a href="/detail" target="_blank" title="${item.title || ''}">${item.title || ''}</a>
                <span class="pull-right">${item.author || '未知作者'}</span>
            `;
            containerEl.appendChild(listItem);
        });
    });
}

// 渲染最新入库区域
function renderNewIn(newInData) {
    const container = document.getElementById('new-in').querySelector('tbody');
    container.innerHTML = '';

    if (!newInData || newInData.length === 0) {
        container.innerHTML = '<tr><td colspan="3" class="text-center">暂无最新入库书籍</td></tr>';
        return;
    }

    newInData.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td class="text-muted hidden-xs">${item.category || ''}</td>
            <td><a href="/detail" target="_blank" title="${item.title || ''}">${item.title || ''}</a></td>
            <td class="text-right fs-12">${item.author || '未知作者'}</td>
        `;
        container.appendChild(row);
    });
}

// 渲染最新更新区域
function renderNewUp(newUpData) {
    const container = document.getElementById('new-up').querySelector('tbody');
    container.innerHTML = '';

    if (!newUpData || newUpData.length === 0) {
        container.innerHTML = '<tr><td colspan="5" class="text-center">暂无最新更新书籍</td></tr>';
        return;
    }

    newUpData.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td class="text-muted hidden-xs">${item.category || ''}</td>
            <td><a href="/detail" target="_blank" title="${item.title || ''}">${item.title || ''}</a></td>
            <td class="hidden-xs"><a href="/detail" target="_blank">${item.latest_chapter || ''}</a></td>
            <td class="text-right fs-12">${item.author || '未知作者'}</td>
            <td class="fs-12 hidden-xs">${item.update_time || ''}</td>
        `;
        container.appendChild(row);
    });
}

// 渲染最新资讯区域
function renderNewMsg(newMsgData) {
    const container = document.getElementById('new-msg').querySelector('tbody');
    container.innerHTML = '';

    if (!newMsgData || newMsgData.length === 0) {
        container.innerHTML = '<tr><td class="text-center">暂无最新资讯</td></tr>';
        return;
    }

    newMsgData.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td><a href="/detail" target="_blank" title="${item.title || ''}">${item.title || ''}</a></td>
        `;
        container.appendChild(row);
    });
}

window.addEventListener('DOMContentLoaded', () => {
    init();

    // 统一处理点击事件
    document.addEventListener('click', (e) => {
        // 排除导航栏中的链接和 More+ 链接
        if (
            (e.target.tagName === 'A' || e.target.tagName === 'BUTTON' || e.target.classList.contains('clickable') || e.target.classList.contains('book-item') || e.target.classList.contains('category-item')) &&
            !e.target.closest('.nav-links li a') &&
            !e.target.classList.contains('pull-right')
        ) {
            e.preventDefault();
            window.location.href = '/detail';
        }
    });
});