import subprocess
import time
import pygame

# Initialize pygame
pygame.init()

# Window settings
WIDTH, HEIGHT = 800, 600
win = pygame.display.set_mode((WIDTH, HEIGHT))
pygame.display.set_caption("Kubernetes Pod Logs Visualization")

# Colors
WHITE = (255, 255, 255)
RED = (255, 0, 0)
GREEN = (0, 255, 0)

# Font settings
font = pygame.font.SysFont(None, 12)

# Get a list of all pods
def get_all_pods():
    cmd = "kubectl get pods -o=jsonpath='{.items[*].metadata.name}'"
    result = subprocess.check_output(cmd, shell=True).decode('utf-8')
    return result.split()

# Get logs of a pod
def get_pod_logs(pod_name):
    cmd = f"kubectl logs {pod_name}"
    try:
        result = subprocess.check_output(cmd, shell=True).decode('utf-8')
        return result
    except:
        return ""

import random

def main():
    clock = pygame.time.Clock()
    pods = get_all_pods()
    prev_logs = {pod: "" for pod in pods}

    RADIUS = 30

    # Generate a random position for each pod inside the window boundaries
    pod_positions = {pod: (random.randint(RADIUS, WIDTH - RADIUS),
                           random.randint(RADIUS, HEIGHT - RADIUS))
                     for pod in pods}

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
                break

        win.fill(WHITE)

        for i, pod in enumerate(pods):
            x, y = pod_positions[pod]

            logs = get_pod_logs(pod)
            color = GREEN if logs != prev_logs[pod] else RED

            pygame.draw.circle(win, color, (x, y), RADIUS)
            label = font.render(pod, 1, (0, 0, 0))
            win.blit(label, (x - (label.get_width() / 2), y - (label.get_height() / 2)))

            prev_logs[pod] = logs

        pygame.display.flip()
        clock.tick(1)  # Update every second

    pygame.quit()

if __name__ == "__main__":
    main()

