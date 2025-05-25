function x2Integrate(xmin, xmax, intervals_count) {
    const dx = (xmax - xmin) / intervals_count;
    let total = 0.0;
    let x = xmin;
    for (let i = 1; i < intervals_count; i++) {
        total = total + dx * (x * x + (x + dx) * (x + dx)) / 2.0;
        x = x + dx;
    }
    return total;
}

module.exports = { x2Integrate };
