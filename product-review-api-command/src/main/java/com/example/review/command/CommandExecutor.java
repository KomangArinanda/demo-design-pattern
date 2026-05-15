package com.example.review.command;

import lombok.RequiredArgsConstructor;
import org.springframework.context.ApplicationContext;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class CommandExecutor {

    private final ApplicationContext applicationContext;

    public <REQUEST, RESPONSE> RESPONSE execute(Class<? extends Command<REQUEST, RESPONSE>> commandClass,
        REQUEST request) {
        Command<REQUEST, RESPONSE> command = applicationContext.getBean(commandClass);
        command.validate(request);
        return command.execute(request);
    }
}
