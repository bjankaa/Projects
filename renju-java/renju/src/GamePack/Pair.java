package GamePack;

// Templat hogy te tudjam tarolni eg valtozoban a poziciot es a sequencet(ebben kersem a mintat)
public class Pair <T1, T2> {
    private final T1 axis ;
    private final T2 value;

    public Pair(T1 axis, T2 value) {
        this.axis=axis;
        this.value=value;
    }
    public T1 getA(){
        return axis;
    }

    public T2 getV(){
        return value;
    }
}
