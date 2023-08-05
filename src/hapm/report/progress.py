"""HAPM Progress reporter"""
from threading import Thread
from time import sleep

STEPS = [
    "▁",
    "▂",
    "▃",
    "▄",
    "▅",
    "▆",
    "▇",
    "█"
]
STEPS_COUNT = len(STEPS)
STEP_DELAY = 0.07

class Progress:
    """Utility class to output Indeterminate progress bar"""
    _title: str
    _running: bool
    _finished: bool

    def __init__(self):
        pass

    def start(self, title: str) -> None:
        """Starts progress with given title"""
        self._title = title
        self._running = True
        self._finished = False
        task = Thread(target=self._show_progress)
        task.daemon = True
        task.start()

    def stop(self) -> None:
        """Stops active progress bar"""
        self._running = False
        while not self._finished:
            sleep(STEP_DELAY)

    def _show_progress(self) -> None:
        states = [
            [0, True],
            [2, True],
            [4, True]
        ]
        prefix =  f"* {self._title} "
        while self._running:
            line = prefix
            for (i, _) in enumerate(states):
                line += STEPS[states[i][0]]
                if states[i][1]:
                    if states[i][0] == STEPS_COUNT - 1:
                        states[i][0] -= 1
                        states[i][1] = False
                    else:
                        states[i][0] += 1
                else:
                    states[i][0] -= 1
                    if states[i][0] == 0:
                        states[i][1] = True
            print(line, end="\r")
            sleep(STEP_DELAY)
        blank = " " * (len(prefix) + STEPS_COUNT)
        print(blank, end="\r")
        self._finished = True
