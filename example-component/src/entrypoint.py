import os

from collect import collect
from act import act
from utils import print_settings_example


print_settings_example()

# The RUN_AS env variable will be passed in so that you know whether to run the
# job as a collector or an actor. This way, you can keep all of your code
# together (some of which will be used for both types of job) and only manage a
# single Dockerfile. If your component is not capable of running a certain type
# of job (i.e. "collector"), then simply fail if it is asked to do so.
RUN_AS = os.getenv('RUN_AS')

if RUN_AS == 'collector':
    collect()
elif RUN_AS == 'actor':
    act()
