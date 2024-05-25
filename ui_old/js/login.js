var jwt = localStorage.getItem("jwt");
if (jwt != null) {
  window.location.href = './index.html'
}

// var expire = localStorage.getItem("expire_on");
// const inputDate = new Date(expire); 

// // Get the current date
// const currentDate = new Date();

// // Compare the input date with the current date
// if (inputDate < currentDate) {
//   console.log('too old jwt token');
//   localStorage.removeItem("jwt");
//   localStorage.removeItem("expire_on");
//   window.location.href = './login.html'
// }

function login() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  const xhttp = new XMLHttpRequest();
//   xhttp.open("POST", "https://www.mecallapi.com/api/login");
    const urlCur = window.location.href;
    const urlAPI = new URL("login", urlCur);
  xhttp.open("POST", urlAPI);
  xhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
  xhttp.send(JSON.stringify({
    "username": username,
    "password": password
  }));
  xhttp.onreadystatechange = function () {
    if (this.readyState == 4) {
      const objects = JSON.parse(this.responseText);
      console.log(objects);
      if (objects['status'] == 'ok') {
        localStorage.setItem("jwt", objects['access_token']);
        localStorage.setItem("expire_on", objects['expire_on']);
        localStorage.setItem("user_id", objects['user_id']);
        localStorage.setItem("user", username);
        Swal.fire({
          text: objects['message'],
          icon: 'success',
          confirmButtonText: 'OK'
        }).then((result) => {
          if (result.isConfirmed) {
            window.location.href = './index.html';
          }
        });
      } else {
        Swal.fire({
          text: objects['message'],
          icon: 'error',
          confirmButtonText: 'OK'
        });
      }
    }
  };
  return false;
}