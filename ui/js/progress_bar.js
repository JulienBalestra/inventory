var MIN_WIDTH = 2.5;


function createProgressBar(current, color, max) {
    var percent_current = (current * 100) / max;
    current = (Math.round(current * 100) / 100);

    var p = '<div class="progress">' +
        '<div class="progress-bar progress-bar-' + color + '" role="progressbar" ' +
        'aria-valuenow="' + current +
        '" aria-valuemin="0" aria-valuemax="100" style="min-width: ' + MIN_WIDTH + 'em; width: ' + percent_current + '%;">' +
        current + '</div></div>';
    return p
}