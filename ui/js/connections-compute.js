function computeMachinesLatencies(machines) {
    var maxLatency = 0;
    var avgLatency = 0;

    for (var i = 0; i < machines.length; i++) {
        var conns = machines[i].Connections;
        if (!conns) {
           continue
        }
        for (var j = 0; j < conns.length; j++) {
            var latency = conns[j].LatencyMs;
            if (latency > maxLatency) {
                maxLatency = latency;
            }
            avgLatency += latency;
        }
    }

    return {maxLatency: maxLatency, avgLatency: avgLatency}
}

function computeConnectionsLatencies(conns) {
    var maxLatency = 0;
    var avgLatency = 0;

    for (var i = 0; i < conns.length; i++) {

        var latency = conns[i].LatencyMs;
        if (latency > maxLatency) {
            maxLatency = latency;
        }
        avgLatency += latency;
    }

    return {maxLatency: maxLatency, avgLatency: avgLatency}
}