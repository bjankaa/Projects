import GamePack.Game;
import GamePack.GameFrame;
import GamePack.WinFrameBlack;
import GamePack.WinFrameWhite;
import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;

public class gameTest {
    @Test
    public void GametestX(){
        GameFrame gf = new GameFrame();
        gf.initGombTomb();
        Game jtk = new Game();
        int x = 4;
        int y= 7;
        String color = "X";
        String[][] gombtomb = gf.gombtomb;
        WinFrameWhite wfw = new WinFrameWhite();
        WinFrameBlack wfb = new WinFrameBlack();
        try {
            jtk.Game(gombtomb, x + 1, y + 1, color, wfw, wfb);
        } catch (IOException e) {
            e.printStackTrace();
        }
        int a= jtk.getStone_x();
        Assert.assertEquals(a,x+1);
    }
    @Test
    public void GametestY(){
        GameFrame gf = new GameFrame();
        gf.initGombTomb();
        Game jtk = new Game();
        int x = 4;
        int y= 7;
        String color = "X";
        String[][] gombtomb = gf.gombtomb;
        WinFrameWhite wfw = new WinFrameWhite();
        WinFrameBlack wfb = new WinFrameBlack();
        try {
            jtk.Game(gombtomb, x + 1, y + 1, color, wfw, wfb);
        } catch (IOException e) {
            e.printStackTrace();
        }
        int b =jtk.getStone_y();
        Assert.assertEquals(b, y+1);
    }
    @Test
    public void GametestColor(){
        GameFrame gf = new GameFrame();
        gf.initGombTomb();
        Game jtk = new Game();
        int x = 4;
        int y= 7;
        String color = "X";
        String[][] gombtomb = gf.gombtomb;
        WinFrameWhite wfw = new WinFrameWhite();
        WinFrameBlack wfb = new WinFrameBlack();
        try {
            jtk.Game(gombtomb, x + 1, y + 1, color, wfw, wfb);
        } catch (IOException e) {
            e.printStackTrace();
        }
        String a= jtk.getColor();
        Assert.assertEquals(a,color);
    }
    @Test
    public void GametestTomb(){
        GameFrame gf = new GameFrame();
        gf.initGombTomb();
        Game jtk = new Game();
        int x = 4;
        int y= 7;
        String color = "X";
        String[][] gombtomb = gf.gombtomb;
        WinFrameWhite wfw = new WinFrameWhite();
        WinFrameBlack wfb = new WinFrameBlack();
        try {
            jtk.Game(gombtomb, x + 1, y + 1, color, wfw, wfb);
        } catch (IOException e) {
            e.printStackTrace();
        }
        Assert.assertNotNull(gombtomb);
    }

}
