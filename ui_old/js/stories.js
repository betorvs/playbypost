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

async function getStories() {
    const userid = localStorage.getItem("user_id");
    const user = localStorage.getItem("user");
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/story" + `/master/` + userid, urlCur);
    const response = await fetch(urlAPI, {
        method: 'GET',
        headers: {
            "X-Username": user,
            "X-Access-Token": jwt }});
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

async function loadStories() {
    const data = await getStories();
    console.log(data);

    const header = `<div class="main-header"><h1>Stories <button type="button" class="button create" onclick="showStoryCreateBox()"><span>Create</span></button></h1></div>`

    const obj = [{ name: "Create Encounter", button: "showEncounterCreateBox"},{ name: "Details", button: "loadStoryDetail"}, { name: "Players", button: "loadPlayersByStoryDetail"}];

    const htmlTable = await formatMain(data, header, obj);

    document.getElementById("content").innerHTML = htmlTable;
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

function storyCreate() {
    const title = document.getElementById("title").value;
    const masterid = document.getElementById("masterid").value;
    const announcement = document.getElementById("announcement").value;
    const notes = document.getElementById("notes").value;
      
    const xhttp = new XMLHttpRequest();
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/story", urlCur);
    xhttp.open("POST", urlAPI);
    xhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhttp.send(JSON.stringify({ 
      "title": title, 
      "master_id": Number(masterid),
      "announcement": announcement, 
      "notes": notes,
    }));
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        const objects = JSON.parse(this.responseText);
        Swal.fire(objects['msg']);
        loadStories();
      }
    };
}

function showEncounterCreateBox(id) {
    Swal.fire({
        title: 'Create encounter',
        html:
            '<input id="title" class="swal2-input" placeholder="Title" max="50">' +
            '<input id="storyid" class="swal2-input" placeholder="Story-ID" type="number"value=' + id + '>' +
            '<input id="announcement" class="swal2-input" placeholder="Announcement">' +
            '<input id="notes" class="swal2-input" placeholder="Notes">',
        focusConfirm: false,
        customClass: {
            confirmButton: "button create",
        },
        preConfirm: () => {
            encounterCreate();
        }
    })
}

async function getStoryByID(id) {
    // const userid = localStorage.getItem("user_id");
    const user = localStorage.getItem("user");
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/story" + `/` + id, urlCur );
    const response = await fetch(urlAPI,{
        method: 'GET',
        headers: {
            "X-Username": user,
            "X-Access-Token": jwt }});
    const data = await response.json();
    return data;
}

async function formatStoryDetail(details, encs, header, suffix) {
    console.log(details);
    console.log(suffix);

    let htmlTable = `<div class="detail-dinamyc">`;
    
    htmlTable += `<h2>Title: ${details.title}</h2>`;
    htmlTable += `<h3>Announcement: ${details.announcement}</h3>`;
    htmlTable += `<p>Notes: ${details.notes}</p><br/>`;
    htmlTable += `</div>`;

    if (!Object.is(encs, null)) {
        htmlTable += `<div class="table-dinamyc" style="overflow-x:auto;"><table>`;

        htmlTable += `<tr>`;
        for (const obj of Object.keys(encs[0])) {
            htmlTable += `<th>${obj}</th>`;
        }
        htmlTable += `<th>Options</th>`;
        htmlTable += `</tr>`;
        for (const d of encs) {
            htmlTable += `<tr>`;
            for (const obj of Object.keys(d)) {
                htmlTable += `<td>${d[obj]}</td>`;
            }
            htmlTable += `<td>`
            for (const s of suffix) {
                htmlTable += `<button type="button" class="button create" onclick="${s.button}(` + details.id + `,` + d["id"] + `)"><span>${s.name}</span></button>`;
            }

            htmlTable += `</td>`
            htmlTable += `</tr>`;
        }
    
        htmlTable += `</table>`;
    }
    htmlTable += `</div>`;
    return header + htmlTable;
}

async function getEncountersByStoryID(id) {
    const user = localStorage.getItem("user");
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/encounter/story" + `/` + id, urlCur);
    const response = await fetch(urlAPI,{
        method: 'GET',
        headers: {
            "X-Username": user,
            "X-Access-Token": jwt }});
    const data = await response.json();
    return data;
}

async function loadStoryDetail(id) {
    const details = await getStoryByID(id);

    const header = `<div class="main-header"><h1>Story Detail </h1></div>`

    const encs = await getEncountersByStoryID(id);
    const obj = [{ name: "Encounter Details", button: "loadEncounterDetail"}];

    const htmlTable = await formatStoryDetail(details, encs, header, obj);

    document.getElementById("content").innerHTML = htmlTable;
}

async function getPlayersByStoryID(id) {
    const user = localStorage.getItem("user");
    const urlCur = window.location.href;
    const urlAPI = new URL("api/v1/player/story" + `/` + id, urlCur);
    const response = await fetch(urlAPI,{
        method: 'GET',
        headers: {
            "X-Username": user,
            "X-Access-Token": jwt }});
    const data = await response.json();
    return data;
}

async function formatPlayersDetail(details, players, header) {

    let htmlTable = `<div class="detail-dinamyc">`;
    
    htmlTable += `<h2>Title: ${details.title}</h2>`;
    htmlTable += `<h3>Announcement: ${details.announcement}</h3>`;
    htmlTable += `<p>Notes: ${details.notes}</p><br/>`;
    htmlTable += `</div>`;

    if (!Object.is(players, null)) {
        // let playersTable = playerTableByStoryID(players);
        // htmlTable += String(playersTable);
        htmlTable += `<div class="table-dinamyc" style="overflow-x:auto;"><table>`;

        htmlTable += `<tr>`;
        htmlTable += `<th>ID</th><th>Name</th><th>Abilities</th><th>Skills</th><th>Extension</th><th>Details</th><th>RPG System</th>`;
        htmlTable += `</tr>`;
        for (const [key, value] of Object.entries(players)) {
            htmlTable += `<tr>`;
            htmlTable += `<td>${key}</td>`;
            htmlTable += `<td>${value.name}</td>`;
            htmlTable += `<td>`;
            for (const [k, v] of Object.entries(value.abilities)) {
                htmlTable += `${k}: ${v}\n`;
            }
            htmlTable += `</td>`;
            htmlTable += `<td>`;
            for (const [k, v] of Object.entries(value.skills)) {
                htmlTable += `${k}: ${v}\n`;
            }
            htmlTable += `</td>`;
            htmlTable += `<td>`;
            for (const [k, v] of Object.entries(value.extension)) {
                htmlTable += `${k}: ${v}\n`;
            }
            htmlTable += `</td>`;
            htmlTable += `<td>`;
            htmlTable += `Destroyed: ${value.destroyed}\n`;
            if (!Object.is(value.details, null)) {
                for (const [k, v] of Object.entries(value.details)) {
                    htmlTable += `${k}: ${v}\n`;
                }
            }
            htmlTable += `</td>`;
            htmlTable += `<td>${value.rpg}</td>`;
            htmlTable += `</tr>`;

        }
        htmlTable += `</table>`;
        htmlTable += `</div>`;
    }
    
    return header + htmlTable;
}

async function loadPlayersByStoryDetail(id) {
    const details = await getStoryByID(id);

    const header = `<div class="main-header"><h1>Players from Story </h1></div>`

    const players = await getPlayersByStoryID(id);
    const obj = [{ name: "Check Health", button: "checkHealthBox"}];

    const htmlTable = await formatPlayersDetail(details, players, header, obj);

    document.getElementById("content").innerHTML = htmlTable;
}
loadStories();