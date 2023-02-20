class colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    GREY = '\033[90m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    ORANGE = '\033[93m'
    RED = '\033[91m'
    END = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def print_temporary(message: str):
    print(f"{colors.GREY}{message}{colors.END}", )
