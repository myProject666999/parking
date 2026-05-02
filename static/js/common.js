var token = localStorage.getItem('token');

function request(url, method, data, successCallback, errorCallback) {
    var ajaxConfig = {
        url: '/api' + url,
        type: method,
        contentType: 'application/json',
        headers: {
            'Authorization': 'Bearer ' + token
        },
        success: function(res) {
            if (res.code === 0) {
                if (successCallback) successCallback(res);
            } else if (res.code === 401) {
                layer.msg('登录已过期，请重新登录', {icon: 2}, function() {
                    localStorage.removeItem('token');
                    localStorage.removeItem('userInfo');
                    window.location.href = '/login';
                });
            } else {
                layer.msg(res.message, {icon: 2});
                if (errorCallback) errorCallback(res);
            }
        },
        error: function(xhr) {
            if (xhr.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('userInfo');
                window.location.href = '/login';
            } else {
                layer.msg('请求失败', {icon: 2});
                if (errorCallback) errorCallback();
            }
        }
    };

    if (data) {
        ajaxConfig.data = JSON.stringify(data);
    }

    $.ajax(ajaxConfig);
}

function get(url, successCallback, errorCallback) {
    request(url, 'GET', null, successCallback, errorCallback);
}

function post(url, data, successCallback, errorCallback) {
    request(url, 'POST', data, successCallback, errorCallback);
}

function put(url, data, successCallback, errorCallback) {
    request(url, 'PUT', data, successCallback, errorCallback);
}

function del(url, successCallback, errorCallback) {
    request(url, 'DELETE', null, successCallback, errorCallback);
}

function formatDate(date, format) {
    if (!date) return '';
    var d = new Date(date);
    var year = d.getFullYear();
    var month = String(d.getMonth() + 1).padStart(2, '0');
    var day = String(d.getDate()).padStart(2, '0');
    var hour = String(d.getHours()).padStart(2, '0');
    var minute = String(d.getMinutes()).padStart(2, '0');
    var second = String(d.getSeconds()).padStart(2, '0');

    if (format === 'date') {
        return year + '-' + month + '-' + day;
    } else if (format === 'time') {
        return hour + ':' + minute + ':' + second;
    } else {
        return year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second;
    }
}

function formatMoney(money) {
    if (money === null || money === undefined) return '0.00';
    return parseFloat(money).toFixed(2);
}

function showConfirm(msg, callback) {
    layer.confirm(msg, {
        icon: 3,
        title: '提示'
    }, function(index) {
        layer.close(index);
        if (callback) callback();
    });
}

function openForm(title, url, width, height) {
    width = width || '600px';
    height = height || '500px';
    
    layer.open({
        type: 2,
        title: title,
        shadeClose: true,
        shade: 0.8,
        area: [width, height],
        content: url
    });
}
