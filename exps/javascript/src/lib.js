function f(x) {
  return x * x;
}

export function x2Integrate(xmin, xmax, intervals_count) {
  const dx = (xmax - xmin) / intervals_count;
  let total = 0;
  let x = xmin;
  for (let i = 0; i < intervals_count; i++) {
    total += dx * (f(x) + f(x + dx)) / 2.0;
    x += dx;
  }
  return total;
}

export function x2IntegrateMock(_xmin, _xmax, _intervals_count) {
  return 0.0;
}

export function fibonacciRecursive(n) {
  if (n <= 1) return n;
  return fibonacciRecursive(n - 1) + fibonacciRecursive(n - 2);
}

export function fibonacciRecursiveMock(_n) {
  return 0;
}

export function fibonacciIterative(n) {
  if (n <= 1) return n;
  let a = 0, b = 1;
  for (let i = 2; i <= n; i++) {
    const temp = b;
    b = a + b;
    a = temp;
  }
  return b;
}

export function fibonacciIterativeMock(_n) {
  return 0;
}

export function multiply(size) {
  let a = 33, b = 10, result = 0;
  for (let i = 0; i < size; i++) {
    result = a * b;
  }
  return result;
}

export function multiplyMock(_size) {
  return 0;
}

export function multiplyVector(size) {
  let a = 33, b = 10;
  const aVector = new Array(size).fill(a);
  const bVector = new Array(size).fill(b);
  const resultVector = new Array(size);
  for (let i = 0; i < size; i++) {
    resultVector[i] = aVector[i] * bVector[i];
  }
  return;
}

export function multiplyVectorMock(_size) {
  return;
}

export function factorize(n) {
  const factors = [];
  let d = 2;
  while (n > 1) {
    while (n % d === 0) {
      factors.push(d);
      n = Math.floor(n / d);
    }
    d++;
  }
  return;
}

export function factorizeMock(_n) {
  return;
}