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

    if (status == "pause") {
        ret.innerHTML = "<span class=\"label label-danger\">Pause</span>";

    } else if (status == "null") {
        ret.innerHTML = "<span class=\"label label-warning\">Waiting</span>";

    } else if (status == "running") {
        ret.innerHTML = "<span class=\"label label-success\">Running</span>";

    } else {
        ret.innerHTML = "<span class=\"label label-danger\">Stop</span>";
    }
}

function compare(a, b) {

    if (a.Hostname > b.Hostname)
        return 1;

    return 0;
}

function machinesTab(m) {
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
            cell.style.width = "80%";
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

function interfacesTab(m) {
    var table = document.getElementById("interfaces_tab");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    insertInterfacesRows(m, table)
}

function insertConnectionsCell(m, row, fields, max) {

    var cell = null;
    var key = null;
    var value = null;


    for (var i = 0; i < fields.length; i++) {

        cell = row.insertCell(i);
        key = fields[i];
        value = m[key];

        if (key == "Connections" && value) {
            var conn = "";
            for (var j = 0; j < value.length; j++) {
                conn += value[j].IPv4 + getProgress(value[j].LatencyMs, max) + "<br>";
            }
            cell.style.width = "70%";
            cell.innerHTML = conn;

        } else {
            cell.innerHTML = value;
        }
    }
}

function getProgress(current, max) {
    var percent_current = (current * 100) / max;
    current = (Math.round(current * 100) / 100);
    var color = "success";
    if (current > 100) {
        color = "danger";
    } else if (current > 30) {
        color = "warning";
    }
    var p = '<div class="progress">' +
        '<div class="progress-bar progress-bar-' + color + '" role="progressbar" ' +
        'aria-valuenow="' + current +
        '" aria-valuemin="0" aria-valuemax="100" style="min-width: 4em; width: ' + percent_current + '%;">' +
        current + ' ms</div></div>';
    return p
}

function sortComputeMax(m) {
    var max = 0;

    for (var i = 0; i < m.length; i++) {
        var conns = m[i].Connections;
        if (conns) {
            conns.sort(function (a, b) {
                if (a.IPv4 > b.IPv4) {
                    return 1;
                }
                return 0;
            });
            for (var j = 0; j < conns.length; j++) {
                if (conns[j].LatencyMs > max) {
                    max = conns[j].LatencyMs;
                }
            }
        }
    }
    return max
}

function insertConnectionsRows(m, table) {

    var fields = ["Hostname", "PublicIP", "Connections"];
    var row = null;

    var skip = 0;

    var max = sortComputeMax(m);

    for (var j = 0; j < m.length; j++) {
        if (m[j].Alive == false) {
            skip++;
            continue
        }
        row = table.insertRow(j - skip);
        insertConnectionsCell(m[j], row, fields, max);
    }

    var index = table.insertRow(0);

    for (var i = 0; i < fields.length; i++) {
        var cell = index.insertCell(i);
        cell.innerHTML = "<b>" + fields[i] + "</b>";
    }
}

function connectionsTab(m) {
    var table = document.getElementById("connections_tab");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    insertConnectionsRows(m, table)
}

function readData(sData) {

    try {
        var m = JSON.parse(sData);
    } catch (e) {
        statusReply("stop");
        return
    }

    if (m == null) {
        statusReply("pause");
        return
    }
    statusReply("running");
    m.sort(compare);

    machinesTab(m);
    interfacesTab(m);
    connectionsTab(m);
}

function request() {
    var xhr = getXMLHttpRequest();

    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 0)) {
            readData(xhr.responseText);
            setTimeout(request, 5000);
        } else if (xhr.readyState == 1 && (xhr.status == 202 || xhr.status == 203)) {
            console.log(xhr.readyState + xhr.status);
            statusReply("stop");

        } else if (xhr.status == 502) {
            statusReply("stop");
        }
    };
    xhr.open("GET", "/api/v0", true);
    xhr.send(null);
}

request();
