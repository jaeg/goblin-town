# goblin-town

# Description
Every evil overlord has to start somewhere, and it looks like you have to start with goblins..
Goblins are selfish creatures that will look out for themselves first.  This creates come challenges when trying to rule over them.

# Controls

1,2,3,4 - Zoom.  3 and 4 are kinda iffy.
W,A,S,D - moves camera around.
Clicking on something on the map will show some information about it.
Clicking on the mini map will move the camera to that location.


# Goblin AI
Goblins can survive on their own and follow a basic set of rules to attempt to do so.
- When idle, if there are too few goblins around and energy is not low seek nearest goblins.
- When idle, if there's enough goblins around and the goblin has enough energy, make a new goblin and give it half of parent's energy.
- If energy is low, take nearest food that the goblin can see.
- If the goblin does not see any food near it the goblin will begin to wander searching for food.
- If attacked, if there are more goblins around me than hostile non-goblins then attack back.  Otherwise run away.

Goblins are social-ish creatures and will want to seek out other goblins in order to have protection in numbers and to reproduce.  They are also greedy creatures and will take whatever food they find nearest to them without consideration of others or how they got it.  If a goblin is attacked either while idle or getting food it will attack back but only if there are more goblins around the attacked goblin than there are hostile non-goblins.  If out numbered a goblin will flee the opposite direction of the attacker.

It's unknown how goblins actually reproduce, but when there's enough of them around they and plenty of food they'll end up producing new goblins in adjacent tiles.