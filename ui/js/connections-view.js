function createCollapseLatency(oneMachine, mLats, cLats) {

    var currentLatency = 0;
    var avgLatency = 0;
    var body = "";

    for (var i = 0; i < oneMachine.Connections.length; i++) {
        var conn = oneMachine.Connections[i];
        currentLatency = conn.LatencyMs;
        avgLatency += currentLatency;
        body += createIpByProgress(conn, cLats);
    }
    avgLatency = avgLatency / oneMachine.Connections.length;

    var header = '<div class="panel panel-default">' +
        '<div class="panel-heading" role="tab" id="headingOne">' +
        '<h4 class="panel-title">' +
        '<a role="button" data-toggle="collapse" data-parent="#accordion" ' +
        'href="#machine' + oneMachine.ID + '"' +
        ' aria-expanded="false" aria-controls="machine"' + oneMachine.ID +'>' +
        createProgressBar(avgLatency, "success", mLats.maxLatency) +
        '</a>' +
        '</h4>' +
        '</div>' +
        '<div id="machine' + oneMachine.ID + '" class="panel-collapse collapse " role="tabpanel" ' +
        'aria-labelledby="machine' + oneMachine.ID + '">' +
        '<div class="panel-body">';

    var footer = '</div></div>';

    return header + body + footer
}

function createIpByProgress(conn, cLats) {

    var content = ' <div class="row">' +
        '<div class="col-md-2">' +
        conn.IPv4 +
        '</div>' +
        '<div class="col-md-10">' +
        createProgressBar(conn.LatencyMs, "success", cLats.maxLatency) +
        '</div>' +
        '</div>';
    return content
}

