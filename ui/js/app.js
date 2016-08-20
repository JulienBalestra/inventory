function getXMLHttpRequest() {

    var xhr = null;

    if (window.XMLHttpRequest || window.ActiveXObject) {
        if (window.ActiveXObject) {
            try {
                xhr = new ActiveXObject("Msxml2.XMLHTTP");
            } catch (e) {
                xhr = new ActiveXObject("Microsoft.XMLHTTP");
            }
        } else {
            xhr = new XMLHttpRequest();
        }
    } else {
        alert("No support for XMLHTTPRequest...");
        return null;
    }
    return xhr;
}

function parseMetadata(meta) {
    if (meta == null) {
        return ""
    }
    if (meta.role == "worker") {
        return meta.role + " " + meta["state"];
    } else {
        return meta.role;
    }
}

function aliveLabel(alive) {
    if (alive) {
        return "<span class=\"label label-success\">True</span>"
    } else {
        return "<span class=\"label label-danger\">False</span>"
    }
}

function insertCells(m, row, order) {

    var cell = null;
    var key = null;
    var value = null;

    for (var i = 0; i < order.length; i++) {
        cell = row.insertCell(i);
        key = order[i];
        value = m[key];

        if (key == "Metadata") {
            cell.innerHTML = parseMetadata(value);
        } else if (key == "Alive") {
            cell.innerHTML = aliveLabel(value)
        } else {
            cell.innerHTML = value;
        }
    }
}

function insertRows(m, table) {

    var order = ["Alive", "ID", "PublicIP", "Hostname", "Metadata"];
    var row = null;

    for (var j = 0; j < m.length; j++) {
        row = table.insertRow(j);
        insertCells(m[j], row, order);
    }

    var index = table.insertRow(0);

    for (var i = 0; i < order.length; i++) {
        var cell = index.insertCell(i);
        cell.innerHTML = order[i];
    }
}

function statusReply(status) {

    var ret = document.getElementById("status");

    if (status == "timeout") {
        ret.innerHTML = "<span class=\"label label-danger\">Timeout</span>";

    } else if (status == "null") {
        ret.innerHTML = "<span class=\"label label-warning\">Empty Reply</span>";

    } else if (status == "healthy") {
        ret.innerHTML = "<span class=\"label label-success\">Healthy</span>";

    } else {
        ret.innerHTML = "<span class=\"label label-danger\">Error</span>";
    }
}

function compare(a, b) {

    if (a.ID < b.ID)
        return 1;

    return 0;
}

function readData(sData) {

    try {
        var m = JSON.parse(sData);
    } catch (e) {
        statusReply("error");
        return
    }

    if (m == null) {
        statusReply("null");
        return
    }
    statusReply("healthy");
    m.sort(compare);

    var table = document.getElementById("machines");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    insertRows(m, table)
}

function request() {
    var xhr = getXMLHttpRequest();

    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 0)) {
            readData(xhr.responseText);
            setTimeout(request, 5000);
        } else if (xhr.readyState == 1 && (xhr.status == 202 || xhr.status == 203)) {
            console.log(xhr.readyState + xhr.status);
            statusReply("timeout");

        } else if (xhr.status == 502) {
            statusReply("error");
        }
    };
    xhr.open("GET", "/api/v0", true);
    xhr.send(null);
}

request();
