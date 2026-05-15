package com.example.review.exception;

public class DuplicateReviewException extends RuntimeException {

    public DuplicateReviewException() {
        super("Duplicate review");
    }
}
