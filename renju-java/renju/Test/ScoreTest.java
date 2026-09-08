import GamePack.Game;
import ScorePack.Score;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class ScoreTest {

    @Test
    public void getLepesszamTest(){
        Score uj = new Score();
        uj.setLepesszam("10");
        uj.setFeherlepes("5");
        uj.setFeketelepes("5");
        uj.setNyertes("fekete");
        String a =uj.getLepesszam();
        Assert.assertEquals("10",a);
    }
    @Test
    public void getFeketeLepesszamTest(){
        Score uj = new Score();
        uj.setLepesszam("10");
        uj.setFeherlepes("5");
        uj.setFeketelepes("5");
        uj.setNyertes("fekete");
        Assert.assertEquals("5",uj.getFeherlepes());
    }
    @Test
    public void getFeherLepesszamTest(){
        Score uj = new Score();
        uj.setLepesszam("10");
        uj.setFeherlepes("5");
        uj.setFeketelepes("5");
        uj.setNyertes("fekete");
        Assert.assertEquals("5",uj.getFeketelepes());
    }
    @Test
    public void getNyertesTest(){
        Score uj = new Score();
        uj.setLepesszam("10");
        uj.setFeherlepes("5");
        uj.setFeketelepes("5");
        uj.setNyertes("fekete");
        Assert.assertEquals("fekete",uj.getNyertes());
    }
}
