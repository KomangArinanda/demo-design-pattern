package com.example.review.command;

public interface Command<REQUEST, RESPONSE> {

    RESPONSE execute(REQUEST request);

    default void validate(REQUEST request) {
    }
}
