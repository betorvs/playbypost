var jwt = localStorage.getItem("jwt");
if (jwt == null) {
  window.location.href = './login.html'
}

var expire = localStorage.getItem("expire_on");
const inputDate = new Date(expire); 

// Get the current date
const currentDate = new Date();

// Compare the input date with the current date
if (inputDate < currentDate) {
  console.log('too old jwt token');
  localStorage.removeItem("jwt");
  localStorage.removeItem("expire_on");
  localStorage.removeItem("user_id");
  localStorage.removeItem("user");
  window.location.href = './login.html'
}

function logout() {
    localStorage.removeItem("jwt");
    localStorage.removeItem("expire_on");
    localStorage.removeItem("user_id");
    localStorage.removeItem("user");
    window.location.href = './login.html'
  }