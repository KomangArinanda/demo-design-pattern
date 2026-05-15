package com.example.review.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class DatabaseLatencySimulator {

    private final long latencyMs;

    public DatabaseLatencySimulator(@Value("${app.database.simulated-latency-ms:0}") long latencyMs) {
        this.latencyMs = latencyMs;
    }

    public void delay() {
        if (latencyMs <= 0) {
            return;
        }
        try {
            Thread.sleep(latencyMs);
        } catch (InterruptedException ex) {
            Thread.currentThread().interrupt();
        }
    }
}
