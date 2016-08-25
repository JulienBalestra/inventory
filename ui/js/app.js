var ALL_LATS = [];
var POLLING = 10000;

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

function getLen(conns) {

    if (conns == null) {
        return 0
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
            cell.innerHTML = getLen(value);
        } else if (key == "Connections") {
            var color = "warning";
            var lat = 0;
            if (value != null) {
                for (var j = 0; j < value.length; j++) {
                    lat += value[j].LatencyMs;
                }
                lat = Math.round((lat / value.length) * 10) / 10;
                if (lat < 20) {
                    color = "success";
                } else if (lat > 50) {
                    color = "danger";
                }

            }
            cell.innerHTML = createProgressBar(lat, color, 20);

        }
        else {
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
            var mask = "";
            for (var j = 0; j < value.length; j++) {
                ifaces += value[j].IPv4 + "</br>";
                mask += value[j].Netmask + "</br>";
            }
            cell.innerHTML = mask;
            cell = row.insertCell(i);
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

    fields.push("Netmask");

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

function insertConnectionsCell(oneMachine, row, mLats, cLats) {

    var i = 0;
    var cell;

    cell = row.insertCell(i++);
    cell.innerHTML = oneMachine.Hostname;

    cell = row.insertCell(i++);
    cell.innerHTML = oneMachine.PublicIP;

    cell = row.insertCell(i);
    cell.style.width = "80%";
    cell.innerHTML = createCollapseLatency(oneMachine, mLats, cLats);

}

function sortComputeLatencies(m) {
    var max = 0;
    var moy = 0;
    var latencies = 0;

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
                moy += conns[j].LatencyMs;
                latencies++;
                if (conns[j].LatencyMs > max) {
                    max = conns[j].LatencyMs;
                }
            }
        }
    }
    moy = moy / latencies;
    return {"moy": moy, "max": max}
}

function insertConnectionsRows(m, mLats, table) {

    var fields = ["Hostname", "PublicIP", "Connections"];
    var row = null;

    var skip = 0;

    for (var j = 0; j < m.length; j++) {
        if (m[j].Alive == false) {
            skip++;
            continue
        }
        row = table.insertRow(j - skip);
        var cLats = computeConnectionsLatencies(m[j].Connections);
        insertConnectionsCell(m[j], row, mLats, cLats);
    }

    var index = table.insertRow(0);

    for (var i = 0; i < fields.length; i++) {
        var cell = index.insertCell(i);
        cell.innerHTML = "<b>" + fields[i] + "</b>";
    }
}


function latencyGraph() {

    var pts = "";
    var x = "";
    var y = "";

    while (ALL_LATS.length > 36) {
        ALL_LATS.shift()
    }

    var maxi = 0;
    var mini = 0;

    for (var i = 0; i < ALL_LATS.length; i++) {
        if (ALL_LATS[i].moy > maxi) {
            maxi = ALL_LATS[i].moy;
        }

        if (ALL_LATS[i].moy < mini) {
            mini = ALL_LATS[i].moy;
        }
    }

    maxi = Math.round(maxi * 1.20);
    mini = Math.round(mini * 0.80);

    if (mini < 0) {
        mini = 0;
    }

    if (mini == 0 && maxi == 0) {
        return
    }

    var padding = 500 / (maxi - mini);
    var ladder = "";

    while (mini < maxi) {
        ladder += '<text x="50" y="' + (500 - (mini * padding)) + '">' + mini + '</text>';
        mini++;
    }
    var nb = 0;
    for (i = 0; i < ALL_LATS.length; i++) {
        x = (nb * 20) + 80;
        y = 500 - (ALL_LATS[i].moy * padding);
        pts += '<circle cx="' + x + '" cy="' + y + '" data-value="' + ALL_LATS[i].moy + '" r="4"></circle>';
        if (i % 5 == 1) {
            // ladder += '<text x="' + (x) + '" y="520">' + Math.round(ALL_LATS.length - (i * (POLLING / 1000))) + '</text>';
        }
        nb++;
    }

    var data = document.getElementById("latency-data");
    var labels = document.getElementById("latency-labels");

    data.innerHTML = pts;
    labels.innerHTML = ladder;
}

function connectionsTab(m) {
    var lats = sortComputeLatencies(m);

    ALL_LATS.push(lats);

    latencyGraph();

    var table = document.getElementById("connections_tab");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    table.deleteRow(0);

    lats = computeMachinesLatencies(m);
    insertConnectionsRows(m, lats, table)
}

function readData(sData) {

    try {
        var m = JSON.parse(sData);
    } catch (e) {
        console.log(e);
        statusReply("stop");
        return
    }

    if (m == null) {
        console.log("empty reply");
        statusReply("stop");
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
            setTimeout(request, POLLING);
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
