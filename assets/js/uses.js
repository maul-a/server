
document.getElementById('logout-button').onclick = function()
{
    deleteCookie('sessionId');
    document.location.href = '/admin'
}