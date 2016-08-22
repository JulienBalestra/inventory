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

function parseInterfaces(interfaces) {

    if (interfaces == null) {
        return ""
    } else {
        return interfaces.length
    }

}

function parseConnections(conns) {

    if (conns == null) {
        return ""
    } else {
        return conns.length
    }

}

function insertMachineCells(m, row, order) {

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
        } else if (key == "Interfaces") {
            cell.innerHTML = parseInterfaces(value);
        } else if (key == "Connections") {
            cell.innerHTML = parseInterfaces(value);
        } else {
            cell.innerHTML = value;
        }
    }
}

function insertMachinesRows(m, table) {

    var fields = ["Alive", "ID", "PublicIP", "Hostname", "Interfaces", "Connections", "Metadata"];
    var row = null;

    for (var j = 0; j < m.length; j++) {
        row = table.insertRow(j);
        insertMachineCells(m[j], row, fields);
    }

    var index = table.insertRow(0);

    for (var i = 0; i < fields.length; i++) {
        var cell = index.insertCell(i);
        cell.innerHTML = "<b>" + fields[i] + "</b>";
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

    if (a.Hostname > b.Hostname)
        return 1;

    return 0;
}

function machines_tab(m) {
    var table = document.getElementById("machines_tab");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    insertMachinesRows(m, table)
}

function insertInterfacesCells(m, row, fields) {

    var cell = null;
    var key = null;
    var value = null;

    for (var i = 0; i < fields.length; i++) {

        cell = row.insertCell(i);
        key = fields[i];
        value = m[key];

        if (key == "Interfaces" && value) {
            var ifaces = "";
            for (var j = 0; j < value.length; j++) {
                ifaces += value[j].IPv4 + " /" + value[j].Netmask + "</br>";
            }
            cell.innerHTML = ifaces;

        } else {
            cell.innerHTML = value;
        }
    }
}

function insertInterfacesRows(m, table) {

    var fields = ["Hostname", "Interfaces"];
    var row = null;


    for (var j = 0; j < m.length; j++) {
        row = table.insertRow(j);
        insertInterfacesCells(m[j], row, fields);
    }

    var index = table.insertRow(0);

    for (var i = 0; i < fields.length; i++) {
        var cell = index.insertCell(i);
        cell.innerHTML = "<b>" + fields[i] + "</b>";
    }
}

function interfaces_tab(m) {
    var table = document.getElementById("interfaces_tab");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    insertInterfacesRows(m, table)
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

    machines_tab(m);
    interfaces_tab(m);
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
