import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;

class DaySixteenFunctions {
	public static void main(String[] args) throws IOException {
		var lines = Files.readAllLines(Path.of("day16/input.txt"));
		System.out.println(partOne(new ArrayList<>(lines)));
		System.out.println(partTwo(new ArrayList<>(lines)));
	}

	public static Long parseInt(String value) {
		Long cur = 0L;
		for (Long i = 0L; (i < Long.valueOf(value.codePointCount(0, value.length()))); i = (i + 1L)) {
			if (((value.codePointAt(value.offsetByCodePoints(0, Long.valueOf(i).intValue())) < 48 /* '0' */) || (value.codePointAt(value.offsetByCodePoints(0, Long.valueOf(i).intValue())) > 57 /* '9' */))) {
				return 0L;
			}
			cur = (cur * 10L);
			cur = (cur + (long)((value.codePointAt(value.offsetByCodePoints(0, Long.valueOf(i).intValue())) - 48 /* '0' */)));
		}
		return cur;
	}
	public static ArrayList<String> split(String s, Integer sep) {
		ArrayList<String> result = (new ArrayList<String>());
		String cur = "";
		for (Long i = 0L; (i < Long.valueOf(s.codePointCount(0, s.length()))); i = (i + 1L)) {
			if (Objects.equals(s.codePointAt(s.offsetByCodePoints(0, Long.valueOf(i).intValue())), sep)) {
				result.add(cur);
				cur = "";
			} else {
				cur = (cur + new String(Character.toChars(s.codePointAt(s.offsetByCodePoints(0, Long.valueOf(i).intValue())))));
			}
		}
		if ((Long.valueOf(cur.codePointCount(0, cur.length())) != 0L)) {
			result.add(cur);
		}
		return result;
	}
	public static String substring(String str, Long start, Long end) {
		String result = "";
		for (Long i = start; ((i < end) && (i < Long.valueOf(str.codePointCount(0, str.length())))); i = (i + 1L)) {
			result = (result + new String(Character.toChars(str.codePointAt(str.offsetByCodePoints(0, Long.valueOf(i).intValue())))));
		}
		return result;
	}
	public static HashSet<String> copySet(HashSet<String> toCopy) {
		HashSet<String> copy = (new HashSet<String>());
		for (String value : toCopy) {
			copy.add(value);
		}
		return copy;
	}
	public static Valve parseValve(String line) {
		ArrayList<String> splitLine = DaySixteenFunctions.split(line, 32 /* ' ' */);
		String name = splitLine.get(Long.valueOf(1L).intValue());
		Long pressure = DaySixteenFunctions.parseInt(DaySixteenFunctions.substring(splitLine.get(Long.valueOf(4L).intValue()), 5L, (Long.valueOf(splitLine.get(Long.valueOf(4L).intValue()).codePointCount(0, splitLine.get(Long.valueOf(4L).intValue()).length())) - 1L)));
		ArrayList<String> tunnels = (new ArrayList<String>());
		for (Long i = 9L; (i < Long.valueOf(splitLine.size())); i = (i + 1L)) {
			String tunnel = splitLine.get(Long.valueOf(i).intValue());
			if ((i != (Long.valueOf(splitLine.size()) - 1L))) {
				tunnel = DaySixteenFunctions.substring(tunnel, 0L, (Long.valueOf(tunnel.codePointCount(0, tunnel.length())) - 1L));
			}
			tunnels.add(tunnel);
		}
		return new Valve(name, pressure, tunnels);
	}
	public static HashMap<String, Valve> getValvesByName(ArrayList<Valve> valves) {
		HashMap<String, Valve> valvesByName = (new HashMap<String, Valve>());
		for (Valve valve : valves) {
			valvesByName.put(valve.name, valve);
		}
		return valvesByName;
	}
	public static ArrayList<State> getNextStates(State curState, HashMap<String, Valve> valvesByName) {
		ArrayList<State> nextStatesForMe = (new ArrayList<State>());
		Long curValvePressure = valvesByName.get(curState.location).pressure;
		if ((!curState.openValves.contains(curState.location) && (curValvePressure != 0L))) {
			State openValveState = curState.copy();
			openValveState.minute = (openValveState.minute + 1L);
			openValveState.releasedPressure = (openValveState.releasedPressure + openValveState.pressurePerMinute);
			openValveState.pressurePerMinute = (openValveState.pressurePerMinute + curValvePressure);
			openValveState.openValves.add(curState.location);
			nextStatesForMe.add(openValveState);
		}
		for (String connection : valvesByName.get(curState.location).tunnels) {
			State moveState = curState.copy();
			moveState.minute = (moveState.minute + 1L);
			moveState.releasedPressure = (moveState.releasedPressure + moveState.pressurePerMinute);
			moveState.location = connection;
			nextStatesForMe.add(moveState);
		}
		ArrayList<State> nextStates = (new ArrayList<State>());
		for (State nextStateForMe : nextStatesForMe) {
			if (Objects.equals(valvesByName.get(nextStateForMe.elephantLocation), null)) {
				nextStates.add(nextStateForMe);
				continue;
			}
			curValvePressure = valvesByName.get(nextStateForMe.elephantLocation).pressure;
			if ((!nextStateForMe.openValves.contains(nextStateForMe.elephantLocation) && (curValvePressure != 0L))) {
				State openValveState = nextStateForMe.copy();
				openValveState.pressurePerMinute = (openValveState.pressurePerMinute + curValvePressure);
				openValveState.openValves.add(nextStateForMe.elephantLocation);
				nextStates.add(openValveState);
			}
			for (String connection : valvesByName.get(nextStateForMe.elephantLocation).tunnels) {
				State moveState = nextStateForMe.copy();
				moveState.elephantLocation = connection;
				nextStates.add(moveState);
			}
		}
		return nextStates;
	}
	public static Long partOne(ArrayList<String> lines) {

		ArrayList<Valve> valves = (new ArrayList<Valve>());
		for (String line : lines) {
			valves.add(DaySixteenFunctions.parseValve(line));
		}
		HashMap<String, Valve> valvesByName = DaySixteenFunctions.getValvesByName(valves);
		Long maxPressure = 0L;
		HashSet<State> curToVisitStates = new HashSet(Set.of(new State(0L, 0L, (new HashSet<String>()), "AA", "nowhere", 0L)));
		HashSet<State> nextToVisitStates = (new HashSet<State>());
		for (; (Long.valueOf(curToVisitStates.size()) > 0L); ) {
			for (State cur : curToVisitStates) {
				if ((cur.releasedPressure > maxPressure)) {
					maxPressure = cur.releasedPressure;
				}
				if ((cur.minute >= 30L)) {
					continue;
				}
				for (State nextState : DaySixteenFunctions.getNextStates(cur, valvesByName)) {
					if (((maxPressure - nextState.releasedPressure) < 100L)) {
						nextToVisitStates.add(nextState);
					}
				}
			}
			curToVisitStates = nextToVisitStates;
			nextToVisitStates = (new HashSet<State>());
		}
		return maxPressure;
	}
	public static Long partTwo(ArrayList<String> lines) {
		ArrayList<Valve> valves = (new ArrayList<Valve>());
		for (String line : lines) {
			valves.add(DaySixteenFunctions.parseValve(line));
		}
		HashMap<String, Valve> valvesByName = DaySixteenFunctions.getValvesByName(valves);
		Long maxPressure = 0L;
		HashSet<State> curToVisitStates = new HashSet(Set.of(new State(0L, 0L, (new HashSet<String>()), "AA", "AA", 0L)));

		HashSet<State> nextToVisitStates = (new HashSet<State>());
		for (; (Long.valueOf(curToVisitStates.size()) > 0L); ) {
			System.out.println(curToVisitStates.iterator().next().minute + " " + curToVisitStates.size());
			for (State cur : curToVisitStates) {
				if ((cur.releasedPressure > maxPressure)) {
					maxPressure = cur.releasedPressure;
				}
				if ((cur.minute >= 26L)) {
					continue;
				}
				for (State nextState : DaySixteenFunctions.getNextStates(cur, valvesByName)) {
					if (((maxPressure - nextState.releasedPressure) < 50)) {
						nextToVisitStates.add(nextState);
					}
				}
			}
			curToVisitStates = nextToVisitStates;
			nextToVisitStates = (new HashSet<State>());
		}
		return maxPressure;
	}
}

class DaySixteenConstants {

}

class Valve {
	Valve(String name, Long pressure, ArrayList<String> tunnels) {
		this.name = name;
		this.pressure = pressure;
		this.tunnels = tunnels;
	}

	String name;
	Long pressure;
	ArrayList<String> tunnels;



}
class State {
	State(Long releasedPressure, Long pressurePerMinute, HashSet<String> openValves, String location, String elephantLocation, Long minute) {
		this.releasedPressure = releasedPressure;
		this.pressurePerMinute = pressurePerMinute;
		this.openValves = openValves;
		this.location = location;
		this.elephantLocation = elephantLocation;
		this.minute = minute;
	}

	Long releasedPressure;
	Long pressurePerMinute;
	HashSet<String> openValves;
	String location;
	String elephantLocation;
	Long minute;

	public State copy() {
		return new State(releasedPressure, pressurePerMinute, DaySixteenFunctions.copySet(openValves), location, elephantLocation, minute);
	}
	@Override
	public boolean equals(Object otherObj) {
		if (otherObj == null) {
			return false;
		}

		if (!(otherObj instanceof State)) {
			return false;
		}

		State other = (State) otherObj;

		return ((((Objects.equals(releasedPressure, other.releasedPressure) && Objects.equals(openValves, other.openValves)) && Objects.equals(location, other.location)) && Objects.equals(elephantLocation, other.elephantLocation)) && Objects.equals(minute, other.minute));
	}
	@Override
	public int hashCode() {
		return Objects.hashCode((new ArrayList(Arrays.asList(Objects.hashCode(releasedPressure), Objects.hashCode(openValves), Objects.hashCode(location), Objects.hashCode(elephantLocation), Objects.hashCode(minute)))));
	}
}