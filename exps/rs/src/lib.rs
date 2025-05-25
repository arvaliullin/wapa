mod utils;

use wasm_bindgen::prelude::*;

fn f(x: f64) -> f64 {
    x * x
}

#[wasm_bindgen]
pub extern "C" fn x2Integrate(xmin: f64, xmax: f64, intervals_count: i32) -> f64 {
    let dx = (xmax - xmin) / (intervals_count as f64);
    let mut total = 0.0;
    let mut x = xmin;
    for _ in 0..intervals_count {
        total += dx * (f(x) + f(x + dx)) / 2.0;
        x += dx;
    }
    total
}

#[wasm_bindgen]
pub extern "C" fn x2IntegrateMock(_xmin: f64, _xmax: f64, _intervals_count: i32) -> f64 {
    0.0
}

#[wasm_bindgen]
pub extern "C" fn fibonacciRecursive(n: i32) -> i32 {
    if n <= 1 {
        n
    } else {
        fibonacciRecursive(n - 1) + fibonacciRecursive(n - 2)
    }
}

#[wasm_bindgen]
pub extern "C" fn fibonacciRecursiveMock(_n: i32) -> i32 {
    0
}

#[wasm_bindgen]
pub extern "C" fn fibonacciIterative(n: i32) -> i32 {
    if n <= 1 {
        return n;
    }
    let mut a = 0;
    let mut b = 1;
    for _ in 2..=n {
        let tmp = a + b;
        a = b;
        b = tmp;
    }
    b
}

#[wasm_bindgen]
pub extern "C" fn fibonacciIterativeMock(_n: i32) -> i32 {
    0
}

#[wasm_bindgen]
pub extern "C" fn multiply(size: i32) -> i32 {
    let a = 33;
    let b = 10;
    let mut result = 0;
    for _ in 0..size {
        result = a * b;
    }
    result
}

#[wasm_bindgen]
pub extern "C" fn multiplyMock(_size: i32) -> i32 {
    0
}

#[wasm_bindgen]
pub extern "C" fn multiplyVector(size: i32) {
    let a = 33;
    let b = 10;
    let a_vector = vec![a; size as usize];
    let b_vector = vec![b; size as usize];
    let mut result_vector = Vec::with_capacity(size as usize);
    for i in 0..(size as usize) {
        result_vector.push(a_vector[i] * b_vector[i]);
    }
}

#[wasm_bindgen]
pub extern "C" fn multiplyVectorMock(_size: i32) {
}

#[wasm_bindgen]
pub extern "C" fn factorize(mut n: i32) {
    let mut d = 2;
    while n > 1 {
        while n % d == 0 {
            n /= d;
        }
        d += 1;
    }
}

#[wasm_bindgen]
pub extern "C" fn factorizeMock(_n: i32) {}
