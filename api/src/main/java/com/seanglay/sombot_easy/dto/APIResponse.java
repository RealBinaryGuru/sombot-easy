package com.seanglay.sombot_easy.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@JsonInclude(JsonInclude.Include.NON_NULL)
@Data
@NoArgsConstructor
@AllArgsConstructor
public class APIResponse<T> {

    private UUID requestId;
    private UUID correlationId;
    private LocalDateTime timestamp;
    private T data;
    private String message;
    private String errorCode;
    private List<String> errors;
    private boolean success;
}
