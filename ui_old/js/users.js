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

async function getUsers() {
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/user", urlCur);
    const response = await fetch(urlAPI);
    const data = await response.json();
    return data;
}

async function formatMain(data, header, suffix) {

    let htmlTable = `<div class="table-dinamyc" style="overflow-x:auto;"><table>`;
    htmlTable += `<tr>`;
    for (const obj of Object.keys(data[0])) {
        htmlTable += `<th>${obj}</th>`;
    }
    htmlTable += `<th>Options</th>`;
    htmlTable += `</tr>`;
    for (const d of data) {
        htmlTable += `<tr>`;
        for (const obj of Object.keys(d)) {
            htmlTable += `<td>${d[obj]}</td>`;
        }
        htmlTable += `<td>`
        for (const s of suffix) {
            htmlTable += `<button type="button" class="button create" onclick="${s.button}(` + d["id"] + `)"><span>${s.name}</span></button>`;
        }
        
        htmlTable += `</td>`
        htmlTable += `</tr>`;
    }
    htmlTable += `</table></div>`;
    return header + htmlTable;
}

async function loadUsers() {
    const data = await getUsers();

    const header = `<div class="main-header"><h1>Users <button type="button" class="button create" onclick="showUserCreateBox()"><span>Create</span></button></h1></div>`

    const obj = [{ name: "Make him a master!", button: "showStoryCreateBox"}];

    const htmlTable = await formatMain(data, header, obj);

    document.getElementById("content").innerHTML = htmlTable;
}

function showUserCreateBox() {
    Swal.fire({
        title: 'Create user',
        html:
            '<input id="username" class="swal2-input" placeholder="Username" max="50">' +
            '<input id="userid" class="swal2-input" placeholder="UserID" max="50">',
        focusConfirm: false,
        customClass: {
            confirmButton: "button create",
        },
        preConfirm: () => {
            userCreate();
        }
    })
}

function userCreate() {
    const username = document.getElementById("username").value;
    const userid = document.getElementById("userid").value;
      
    const xhttp = new XMLHttpRequest();
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/user", urlCur);
    xhttp.open("POST", urlAPI);
    xhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhttp.send(JSON.stringify({ 
      "username": username, "user_id": userid,
    }));
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        const objects = JSON.parse(this.responseText);
        Swal.fire(objects['msg']);
        loadUsers();
      }
    };
}

function showStoryCreateBox(id) {
    Swal.fire({
        title: 'Create story',
        html:
            '<label for="swal2-input" class="swal2-input-label">Title</label>' +
            '<input id="title" class="swal2-input" placeholder="Title" max="50">' +
            '<label for="swal2-input" class="swal2-input-label">Master ID</label>' +
            '<input id="masterid" class="swal2-input" placeholder="MasterID" type="number" value=' + id + '>' +
            '<label for="swal2-textarea" class="swal2-input-label">Announcement</label>' +
            '<textarea id="announcement" class="swal2-textarea" style="display; flex;" placeholder="Announcement when starting it"></textarea>' +
            '<label for="swal2-textarea" class="swal2-input-label">Notes</label>' +
            '<textarea id="notes" class="swal2-textarea" style="display; flex;" placeholder="Notes for your use only"></textarea>',
        focusConfirm: false,
        customClass: {
            confirmButton: "button create",
        },
        preConfirm: () => {
            storyCreate();
        }
    })
}

loadUsers();