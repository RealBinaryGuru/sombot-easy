package com.seanglay.sombot_easy.controller;

import com.seanglay.sombot_easy.dto.APIResponse;
import com.seanglay.sombot_easy.exception.ApplicationAbstractException;
import com.seanglay.sombot_easy.exception.TraceableApplicationAbstractException;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.UUID;

public class BaseController {

    protected <T> ResponseEntity<APIResponse<T>> ok(T data) {
        APIResponse<T> response = new APIResponse<>(UUID.randomUUID(), UUID.randomUUID(), LocalDateTime.now(), data, null, null, null, true);
        return new ResponseEntity<>(response, HttpStatus.OK);
    }

    protected <T> ResponseEntity<APIResponse<T>> ok(T data, UUID requestId, UUID correlationId) {
        APIResponse<T> response = new APIResponse<>(requestId, correlationId, LocalDateTime.now(), data, null, null, null, true);
        return new ResponseEntity<>(response, HttpStatus.OK);
    }

    protected <T> ResponseEntity<APIResponse<T>> success(T data, HttpStatus status) {
        APIResponse<T> response = new APIResponse<>(UUID.randomUUID(), UUID.randomUUID(), LocalDateTime.now(), data, null, null, null, true);
        return new ResponseEntity<>(response, status);
    }

    protected ResponseEntity<APIResponse<?>> error(Exception ex) {
        APIResponse<Object> response = new APIResponse<>(UUID.randomUUID(), UUID.randomUUID(), LocalDateTime.now(), null, "There was an error processing the request", null, Collections.singletonList(ex.getMessage()), false);
        return new ResponseEntity<>(response, HttpStatus.INTERNAL_SERVER_ERROR);
    }

    protected ResponseEntity<APIResponse<?>> error(Exception ex, int statusErrorCode) {
        APIResponse<Object> response = new APIResponse<>(UUID.randomUUID(), UUID.randomUUID(), LocalDateTime.now(), null, "There was an error processing the request", null, Collections.singletonList(ex.getMessage()), false);
        return new ResponseEntity<>(response, HttpStatusCode.valueOf(statusErrorCode));
    }

    protected ResponseEntity<APIResponse<?>> error(ApplicationAbstractException ex) {
        APIResponse<Object> response = new APIResponse<>(UUID.randomUUID(), UUID.randomUUID(), LocalDateTime.now(), null, null, ex.getErrorCode(), Collections.singletonList(ex.getMessage()), false);
        return new ResponseEntity<>(response, HttpStatusCode.valueOf(ex.getStatusCode()));
    }

    protected ResponseEntity<APIResponse<?>> error(TraceableApplicationAbstractException ex) {
        APIResponse<Object> response = new APIResponse<>(ex.getRequestId(), ex.getCorrelationId(), ex.getTimestamp(), null, null, ex.getErrorCode(), Collections.singletonList(ex.getMessage()), false);
        return new ResponseEntity<>(response, HttpStatusCode.valueOf(ex.getStatusCode()));
    }
}
