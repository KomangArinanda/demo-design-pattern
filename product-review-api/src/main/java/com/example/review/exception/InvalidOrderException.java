package com.example.review.exception;

public class InvalidOrderException extends RuntimeException {

    public InvalidOrderException() {
        super("Customer did not purchase this product");
    }
}
