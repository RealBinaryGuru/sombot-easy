package com.seanglay.sombot_easy.exception;

import lombok.Getter;

@Getter
public abstract class ApplicationAbstractException extends Exception {

    private final int statusCode;
    private String errorCode;

    public ApplicationAbstractException(String message, int statusCode) {
        super(message);
        this.statusCode = statusCode;
    }

    public ApplicationAbstractException(String message, int statusCode, String errorCode) {
        this(message, statusCode);
        this.errorCode = errorCode;
    }

}
