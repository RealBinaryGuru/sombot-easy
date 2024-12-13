package com.seanglay.sombot_easy.exception;

import java.time.LocalDateTime;
import java.util.UUID;

public class ApplicationAPIException extends TraceableApplicationAbstractException {

    public ApplicationAPIException(String message, int statusCode) {
        super(message, statusCode);
    }

    public ApplicationAPIException(String message, int statusCode, String errorCode) {
        super(message, statusCode, errorCode);
    }

    public ApplicationAPIException(String message, int statusCode, String errorCode, UUID requestId, UUID correlationId) {
        super(message, statusCode, errorCode, requestId, correlationId);
    }

    public ApplicationAPIException(String message, int statusCode, String errorCode, UUID requestId, UUID correlationId, LocalDateTime timestamp) {
        super(message, statusCode, errorCode, requestId, correlationId, timestamp);
    }
}
