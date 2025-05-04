function f(x) {
  return x * x;
}

export function x2Integrate(xmin, xmax, intervals_count) {
  const dx = (xmax - xmin) / intervals_count;
  let total = 0;
  let x = xmin;

  for (let i = 1; i < intervals_count; i++) {
    total += dx * (f(x) + f(x + dx)) / 2.0;
    x += dx;
  }

  return total;
}

export function x2IntegrateMock(_xmin, _xmax, _intervals_count) {
  return 0.0;
}
