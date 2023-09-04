package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import java.util.*;


@SpringBootApplication
@RestController
public class DemoApplication {
    private static Map<String, String> envVars;
    private static List<String> remoteAddrs;
    private static int sleepTime;

    private static void parseConfigs() {
        envVars = System.getenv();
        remoteAddrs = new ArrayList<String>();
        if (envVars.get("remote_addrs") != null) {
            Collections.addAll(remoteAddrs, envVars.get("remote_addrs").split(","));
        }
        String sleepTimeEnv = envVars.get("sleep_time");
        if (sleepTimeEnv != null) {
            sleepTime = Integer.parseInt(sleepTimeEnv);
        } else {
            sleepTime = 0;
        }
        System.out.println(remoteAddrs);
        System.out.println(sleepTime);
    }

    public static void main(String[] args) {
        parseConfigs();
        SpringApplication.run(DemoApplication.class, args);
    }

    @GetMapping("/*")
    public String sayHello() {
        try {
            Thread.sleep(sleepTime);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        if (remoteAddrs.size() > 0) {
            for (String remoteAddr: remoteAddrs) {
                RestTemplate restTemplate = new RestTemplate();
                restTemplate.getForEntity(remoteAddr, String.class);
                System.out.printf("Call service %s done\n", remoteAddr);
            }
            return String.format("Hello, I have called %s service\n", remoteAddrs.size());
        } else {
            return String.format("Hello, I'm the final service!\n");
        }
    }

}