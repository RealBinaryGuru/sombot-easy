package com.seanglay.sombot_easy.dto;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class APIResponse<T> {

    private UUID requestId;
    private UUID correlationId;
    private LocalDateTime timestamp;
    private T data;
    private String message;
    private String errorCode;
    private List<String> errors;
    private boolean success;

    public APIResponse() {
    }

    public APIResponse(UUID requestId, UUID correlationId, LocalDateTime timestamp, T data, String message, String errorCode, List<String> errors, boolean success) {
        this.requestId = requestId;
        this.correlationId = correlationId;
        this.timestamp = timestamp;
        this.data = data;
        this.message = message;
        this.errorCode = errorCode;
        this.errors = errors;
        this.success = success;
    }

    public UUID getRequestId() {
        return requestId;
    }

    public void setRequestId(UUID requestId) {
        this.requestId = requestId;
    }

    public UUID getCorrelationId() {
        return correlationId;
    }

    public void setCorrelationId(UUID correlationId) {
        this.correlationId = correlationId;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(LocalDateTime timestamp) {
        this.timestamp = timestamp;
    }

    public T getData() {
        return data;
    }

    public void setData(T data) {
        this.data = data;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public void setErrorCode(String errorCode) {
        this.errorCode = errorCode;
    }

    public List<String> getErrors() {
        return errors;
    }

    public void setErrors(List<String> errors) {
        this.errors = errors;
    }

    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }
}
